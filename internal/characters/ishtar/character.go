package ishtar

import (
	"context"
	"time"

	"github.com/Sinketsu/artifactsmmo/internal/events"
	"github.com/Sinketsu/artifactsmmo/internal/generic"
	"github.com/Sinketsu/artifactsmmo/internal/strategy"
)

type Character struct {
	generic.Character

	what     string
	strategy strategy.Strategy
}

func NewCharacter(params generic.Params) *Character {
	gc, err := generic.NewCharacter(params)
	if err != nil {
		panic(err)
	}

	return &Character{
		Character: *gc,
	}
}

func (c *Character) Live(ctx context.Context, events *events.Service) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := c.do(events)
			if err != nil {
				c.Log(err)
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (c *Character) do(events *events.Service) error {
	switch {
	case c.Data().MiningLevel < 35:
		c.setStrategy(
			"mine gold and craft gold bars",
			strategy.NewSimpleGatherStrategy().
				Gather("gold_rocks").
				Craft("gold").
				Bank("gold", "sapphire", "ruby", "emerald", "topaz"),
		)
	case c.Data().WoodcuttingLevel < 35:
		c.setStrategy(
			"gather dead wood and craft dead wood planks (allow some events)",
			strategy.NewSimpleGatherStrategy().
				AllowEvents(events, "Strange Apparition").
				Gather("dead_tree").
				Craft("dead_wood_plank").
				Sell("shrimp", "iron_ore", "spruce_wood", "yellow_slimeball", "red_slimeball", "gold_ore").
				Bank("dead_wood_plank", "sap", "diamond", "strange_ore"),
		)
	case c.Data().FishingLevel < 40:
		c.setStrategy(
			"fishing bass and sell it (allow all events)",
			strategy.NewSimpleGatherStrategy().
				AllowEvents(events, "Strange Apparition", "Magic Apparition").
				Gather("bass_fishing_spot").
				Craft("dead_wood_plank").
				Sell("bass").
				Bank("dead_wood_plank", "sap", "diamond", "strange_ore", "magic_wood", "magic_sap"),
		)
	default:
		c.setStrategy(
			"nothing to do",
			strategy.EmptyStrategy(),
		)
	}

	return c.strategy.Do(&c.Character)
}

func (c *Character) setStrategy(what string, newStrategy strategy.Strategy) {
	if c.what != what {
		c.Log("change strategy:", what)
		c.strategy = newStrategy
		c.what = what
	}
}
