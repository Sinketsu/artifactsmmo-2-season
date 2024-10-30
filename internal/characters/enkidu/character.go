package enkidu

import (
	"context"
	"strings"
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
	gc, err := generic.NewCharacter(client, generic.Params{Name: "Enkidu"}, bank, events)
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
	items := []string{}

	// if c.Data().GearcraftingLevel >= 20 && c.Data().GearcraftingLevel <= 30 {
	// 	items = append(items, "magic_wizard_hat", "steel_helm", "steel_boots", "steel_armor", "steel_legs_armor",
	// 		"skeleton_pants", "skeleton_armor", "skeleton_helmet", "serpent_skin_legs_armor", "steel_shield",
	// 		"tromatising_mask", "serpent_skin_armor")
	// }

	// if c.Data().GearcraftingLevel >= 25 && c.Data().GearcraftingLevel <= 30 {
	// 	items = append(items, "lizard_skin_armor", "lizard_skin_legs_armor", "piggy_pants")
	// }

	if c.Data().JewelrycraftingLevel >= 20 && c.Data().JewelrycraftingLevel <= 30 {
		items = append(items, "ring_of_chance", "dreadful_ring", "steel_ring", "skull_ring", "dreadful_amulet",
			"skull_amulet")
	}

	c.setStrategy(
		"craft something of: "+strings.Join(items, ", "),
		strategy.NewSimpleCraftStrategy().
			WithdrawGold().
			Buy(map[string]int{
				"iron":           500,
				"wolf_bone":      600,
				"skeleton_bone":  500,
				"serpent_skin":   500,
				"cowhide":        500,
				"hardwood_plank": 500,
				"flying_wing":    500,
			}).
			Craft(items...).
			Recycle(items...),
	)

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
