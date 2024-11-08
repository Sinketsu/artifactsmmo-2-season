package ereshkigal

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
	gc, err := generic.NewCharacter(client, generic.Params{Name: "Ereshkigal"}, bank, events)
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
		"do monster tasks",
		strategy.NewTasksFightStrategy().
			Deposit("owlbear_hair", "red_cloth", "skeleton_bone",
				"vampire_blood", "ogre_eye", "ogre_skin",
				"demon_horn", "piece_of_obsidian", "magic_stone", "cursed_book",
				"demoniac_dust", "piece_of_obsidian", "lizard_skin", "tasks_coin").
			DepositGold().
			Sell("mushroom", "red_slimeball", "yellow_slimeball", "blue_slimeball", "green_slimeball",
				"raw_beef", "milk_bucket", "cowhide", "raw_wolf_meat", "wolf_bone", "wolf_hair",
				"raw_chicken", "egg", "feather", "pig_skin", "flying_wing", "skeleton_skull",
				"serpent_skin", "bandit_armor", "golden_egg").
			CancelTasks("lich", "bat", "cultist_acolyte").
			AllowEvents("Bandit Camp", "Portal"),
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
