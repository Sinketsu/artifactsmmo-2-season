package fighter

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
			if slot.Code == "raw_chicken" {
				q = slot.Quantity
			}
		}
		if q == 0 {
			return fmt.Errorf("unexpected")
		}

		geItem, err := c.GetGEItem("raw_chicken")
		if err != nil {
			return fmt.Errorf("get ge item: %w", err)
		}

		q = int(math.Min(float64(q), float64(geItem.MaxQuantity)))
		gold, err := c.Sell("raw_chicken", q, geItem.SellPrice.Value)
		if err != nil {
			return fmt.Errorf("sell: %w", err)
		}

		fmt.Println("sold", q, "raw_chicken", "Earned", gold, "gold")

		// eggs
		q = 0
		for _, slot := range c.Data().Inventory {
			if slot.Code == "egg" {
				q = slot.Quantity
			}
		}
		if q == 0 {
			return fmt.Errorf("unexpected")
		}

		geItem, err = c.GetGEItem("egg")
		if err != nil {
			return fmt.Errorf("get ge item: %w", err)
		}

		q = int(math.Min(float64(q), float64(geItem.MaxQuantity)))
		gold, err = c.Sell("egg", q, geItem.SellPrice.Value)
		if err != nil {
			return fmt.Errorf("sell: %w", err)
		}

		fmt.Println("sold", q, "egg", "Earned", gold, "gold")

		// feather
		q = 0
		for _, slot := range c.Data().Inventory {
			if slot.Code == "raw_chicken" {
				q = slot.Quantity
			}
		}
		if q == 0 {
			return fmt.Errorf("unexpected")
		}

		geItem, err = c.GetGEItem("feather")
		if err != nil {
			return fmt.Errorf("get ge item: %w", err)
		}

		if q <= 30 {
			return nil
		}

		q -= 30

		gold, err = c.Sell("feather", q, geItem.SellPrice.Value)
		if err != nil {
			return fmt.Errorf("sell: %w", err)
		}

		fmt.Println("sold", q, "feather", "Earned", gold, "gold")

		return nil
	}

	err := c.Move(3, -2) // green slime
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	fight, err := c.Fight()
	if err != nil {
		return err
	}

	fmt.Println("got", fight.Xp, "XP")
	fmt.Println("got", fight.Gold, "gold")
	for _, item := range fight.Drops {
		fmt.Println("got", item.Quantity, item.Code)
	}

	return nil
}
