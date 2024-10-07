package ishtar

import (
	"context"
	"time"

	"github.com/Sinketsu/artifactsmmo/internal/generic"
	"github.com/Sinketsu/artifactsmmo/internal/strategy"
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
	if c.Data().FishingLevel < 20 {
		return strategy.NewSimpleGatherStrategy().
			Gather("shrimp_fishing_spot").
			Bank("shrimp", "golden_shrimp").
			Do(&c.Character)
	}

	if c.Data().MiningLevel < 30 {
		return strategy.NewSimpleGatherStrategy().
			Gather("coal_rocks").
			Bank("coal", "topaz", "emerald", "ruby", "sapphire").
			Do(&c.Character)
	}

	if c.Data().WoodcuttingLevel < 30 {
		return strategy.NewSimpleGatherStrategy().
			Gather("birch_tree").
			Bank("coal", "topaz", "emerald", "ruby", "sapphire", "sap", "birch_wood").
			Do(&c.Character)
	}

	if c.Data().FishingLevel < 30 {
		return strategy.NewSimpleGatherStrategy().
			Gather("trout_fishing_spot").
			Bank("sap", "birch_wood", "trout").
			Do(&c.Character)
	}

	return nil
}
