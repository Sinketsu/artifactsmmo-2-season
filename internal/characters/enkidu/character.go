package enkidu

import (
	"context"
	"time"

	"github.com/Sinketsu/artifactsmmo/internal/generic"
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
	if c.Data().WeaponcraftingLevel < 10 {
		return generic.NewSimpleCraftStrategy().Craft("copper_dagger").Recycle("copper_dagger").Do(&c.Character)
	}

	if c.Data().JewelrycraftingLevel < 10 {
		return generic.NewSimpleCraftStrategy().Craft("copper_ring").Recycle("copper_ring").Do(&c.Character)
	}

	if c.Data().GearcraftingLevel < 10 {
		return generic.NewSimpleCraftStrategy().Craft("copper_helmet").Recycle("copper_helmet").Do(&c.Character)
	}

	time.Sleep(1 * time.Second)
	c.Log("idle...")
	return nil
}
