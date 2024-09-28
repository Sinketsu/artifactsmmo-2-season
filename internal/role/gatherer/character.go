package gatherer

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/Sinketsu/artifactsmmo/internal/role/generic"
)

type Character struct {
	generic.Character
}

func NewCharacter(params generic.Params) (*Character, error) {
	gc, err := generic.NewCharacter(params)
	if err != nil {
		return nil, err
	}

	return &Character{
		Character: *gc,
	}, nil
}

func (c *Character) Live(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := c.do()
			if err != nil {
				fmt.Println(err)
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (c *Character) do() error {
	if len(c.Data().Inventory) == c.Data().InventoryMaxItems {
		err := c.Move(5, 1) // Grand Exchange
		if err != nil {
			return err
		}

		q := 0
		for _, slot := range c.Data().Inventory {
			if slot.Code == "copper_ore" {
				q = slot.Quantity
			}
		}
		if q == 0 {
			return fmt.Errorf("unexpected")
		}

		geItem, err := c.GetGEItem("copper_ore")
		if err != nil {
			return err
		}

		gold, err := c.Sell("copper_ore", int(math.Min(float64(q), float64(geItem.MaxQuantity))), geItem.SellPrice.Value)
		if err != nil {
			return err
		}

		fmt.Println("sold", q, "copper_ore.", "Earned", gold, "gold")

		err = c.Move(2, 1) // copper_ore
		if err != nil {
			return err
		}

		return nil
	}

	drop, err := c.Gather()
	if err != nil {
		return err
	}

	fmt.Println("got", drop.Xp, "XP")
	for _, item := range drop.Items {
		fmt.Println("got", item.Quantity, item.Code)
	}

	return nil
}
