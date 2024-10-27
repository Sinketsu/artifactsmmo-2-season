package cetcalcoatl

import (
	"context"
	"time"

	"github.com/Sinketsu/artifactsmmo/internal/api"
	"github.com/Sinketsu/artifactsmmo/internal/generic"
	"github.com/Sinketsu/artifactsmmo/internal/strategy"
)

type Character struct {
	generic.Character

	what     string
	strategy strategy.Strategy
}

func NewCharacter(client *api.Client, bank generic.Bank, events generic.Events, logGroup string, logToken string) *Character {
	gc, err := generic.NewCharacter(client, generic.Params{Name: "Cetcalcoatl"}, bank, events, generic.LogOptions{Group: logGroup, Token: logToken})
	if err != nil {
		panic(err)
	}

	return &Character{
		Character: *gc,
	}
}

func (c *Character) Live(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := c.do()
			if err != nil {
				c.Logger().Error(err.Error())
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (c *Character) do() error {
	c.setStrategy(
		"gather coal_rocks",
		strategy.NewSimpleGatherStrategy().
			AllowEvents("Strange Apparition", "Magic Apparition").
			Gather("coal_rocks").
			Deposit("coal", "topaz", "emerald", "ruby", "sapphire", "strange_ore", "diamond", "magic_wood", "magic_sap").
			DepositGold(),
	)

	// c.setStrategy("craft steel", strategy.NewSimpleCraftStrategy().Craft("steel").Bank("steel"))

	// c.setStrategy("player control", strategy.EmptyStrategy())

	return c.strategy.Do(&c.Character)
}

func (c *Character) setStrategy(what string, newStrategy strategy.Strategy) {
	if c.what != what {
		c.Logger().Info("change strategy: " + what)
		c.strategy = newStrategy
		c.what = what
	}
}
