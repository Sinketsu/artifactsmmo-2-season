package crafter

import (
	"context"
	"fmt"
	"time"

	api "github.com/Sinketsu/artifactsmmo/gen/oas"
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
	if c.InInventory("copper") >= 8 {
		if err := c.Move(2, 1); err != nil {
			return fmt.Errorf("move: %w", err)
		}

		var result api.SkillDataSchemaDetails
		var err error

		switch {
		case c.Data().WeaponcraftingLevel < 5:
			result, err = c.Craft("copper_dagger", 1)
			if err != nil {
				return fmt.Errorf("craft: %w", err)
			}
		case c.Data().JewelrycraftingLevel < 5:
			result, err = c.Craft("copper_ring", 1)
			if err != nil {
				return fmt.Errorf("craft: %w", err)
			}
		case c.Data().GearcraftingLevel < 5:
			result, err = c.Craft("copper_helmet", 1)
			if err != nil {
				return fmt.Errorf("craft: %w", err)
			}
		default:
			return fmt.Errorf("all skills are 5 level")
		}

		c.Log("got", result.Xp, "by crafting", result.Items[0].Code)
		return nil
	}

	if c.InInventory("copper_ore") >= 8 {
		if err := c.Move(1, 5); err != nil {
			return fmt.Errorf("move: %w", err)
		}

		result, err := c.Craft("copper", 1)
		if err != nil {
			return fmt.Errorf("craft: %w", err)
		}

		c.Log("got", result.Xp, "by crafting", result.Items[0].Code)
		return nil
	}

	if q, _ := c.InBank("copper"); q >= 8 && c.Data().InventoryMaxItems-c.InventoryItemCount() >= q {
		if err := c.Move(4, 1); err != nil {
			return fmt.Errorf("move: %w", err)
		}

		err := c.Withdraw("copper", q)
		if err != nil {
			return fmt.Errorf("withdraw: %w", err)
		}

		c.Log("withdraw", q, "copper")
		return nil
	}

	return c.MacroGather("copper_rocks")
}
