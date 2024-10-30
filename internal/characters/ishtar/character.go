package ishtar

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

func NewCharacter(client *api.Client, bank generic.Bank, events generic.Events) *Character {
	gc, err := generic.NewCharacter(client, generic.Params{Name: "Ishtar"}, bank, events)
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
	// c.setStrategy(
	// 	"gather iron_rocks",
	// 	strategy.NewSimpleGatherStrategy().
	// 		AllowEvents("Strange Apparition", "Magic Apparition").
	// 		Gather("iron_rocks").
	// 		Craft("iron").
	// 		DepositGold().
	// 		Deposit("iron", "topaz", "emerald", "ruby", "sapphire", "strange_ore", "diamond", "magic_wood", "magic_sap"),
	// )
	c.setStrategy(
		"fight cow",
		strategy.NewSimpleFightStrategy().
			Fight("cow").
			Deposit("skeleton_bone", "cowhide").
			Sell("skeleton_skull", "raw_beef", "milk_bucket").
			DepositGold(),
	)

	// c.setStrategy("craft steel", strategy.NewSimpleCraftStrategy().Craft("steel").Deposit("steel").DepositGold())

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
