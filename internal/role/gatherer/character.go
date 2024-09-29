package gatherer

import (
	"context"
	"fmt"
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
				c.Log(err)
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (c *Character) do() error {
	if c.Data().FishingLevel < 10 {
		return c.MacroGather("gudgeon", "gudgeon")
	}
	if c.Data().WoodcuttingLevel < 10 {
		return c.MacroGather("ash_tree", "gudgeon", "ash_wood")
	}

	if c.InInventory("copper_ore") > 24 {
		if q, _ := c.InBank("copper_ore"); q == 0 {
			if err := c.Move(4, 1); err != nil {
				return fmt.Errorf("move: %w", err)
			}

			if err := c.Deposit("copper_ore", 24); err != nil {
				return fmt.Errorf("deposit: %w", err)
			}

			return nil
		}
	}

	return c.MacroGather("copper_rocks")
}
