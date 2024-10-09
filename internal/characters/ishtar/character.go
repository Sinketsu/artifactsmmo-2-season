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
	if c.Data().WoodcuttingLevel < 35 {
		return strategy.NewSimpleGatherStrategy().
			AllowEvents(events, "Strange Apparition").
			Gather("dead_tree").
			Craft("dead_wood_plank").
			Bank("srimp", "trout", "dead_wood_plank", "sap", "diamond", "strange_ore").
			Do(&c.Character)
	}

	return strategy.NewSimpleGatherStrategy().
		AllowEvents(events, "Strange Apparition", "Magic Apparition").
		Gather("bass_fishing_spot").
		Bank("srimp", "bass", "trout", "dead_wood_plank", "sap", "diamond", "strange_ore", "magic_wood", "magic_sap").
		Do(&c.Character)
}
