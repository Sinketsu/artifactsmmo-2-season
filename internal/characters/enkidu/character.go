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

	items = append(items, "gold_platelegs", "gold_mask", "gold_helm")
	items = append(items, "gold_ring")

	c.setStrategy(
		"craft something of: "+strings.Join(items, ", "),
		strategy.NewSimpleCraftStrategy().
			WithdrawGold().
			Buy(map[string]int{
				"lizard_skin":    2000,
				"red_cloth":      2000,
				"vampire_blood":  2000,
				"ogre_skin":      1000,
				"demon_horn":     3000,
				"skeleton_skull": 1000,
				"owlbear_hair":   3000,
				"wolf_bone":      2000,
				"skeleton_bone":  2000,
			}).
			Craft(items...).
			Recycle(items...),
	)

	return c.strategy.Do(&c.Character)
}

func (c *Character) setStrategy(what string, newStrategy strategy.Strategy) {
	if c.what != what {
		c.Logger().Info("change strategy: " + what)
		c.strategy = newStrategy
		c.what = what
	}
}
