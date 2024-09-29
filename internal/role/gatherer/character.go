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
	if c.InventoryItemCount() == c.Data().InventoryMaxItems {
		err := c.Move(5, 1) // Grand Exchange
		if err != nil {
			return fmt.Errorf("move: %w", err)
		}

		q := 0
		for _, slot := range c.Data().Inventory {
			if slot.Code == "gudgeon" {
				q = slot.Quantity
			}
		}
		if q == 0 {
			return fmt.Errorf("unexpected")
		}

		geItem, err := c.GetGEItem("gudgeon")
		if err != nil {
			return fmt.Errorf("get ge item: %w", err)
		}

		q = int(math.Min(float64(q), float64(geItem.MaxQuantity)))
		gold, err := c.Sell("gudgeon", q, geItem.SellPrice.Value)
		if err != nil {
			return fmt.Errorf("sell: %w", err)
		}

		fmt.Println("sold", q, "gudgeon", "Earned", gold, "gold")

		return nil
	}

	err := c.Move(4, 2) // fishing
	if err != nil {
		return fmt.Errorf("move: %w", err)
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
