package ereshkigal

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
	return strategy.NewSimpleFightStrategy().
		Bank("owlbear_hair", "golden_egg", "red_cloth", "skeleton_bone", "skeleton_skull",
			"vampire_blood", "flying_wing", "serpent_skin", "ogre_eye", "ogre_skin",
			"bandit_armor", "demon_horn", "piece_of_obsidian").
		Sell("mushroom", "red_slimeball", "yellow_slimeball", "blue_slimeball", "green_slimeball",
			"raw_beef", "milk_bucket", "cowhide", "raw_wolf_meat", "wolf_bone", "wolf_hair",
			"raw_chicken", "egg", "feather", "pig_skin", "lizard_skin").
		CancelTasks("lich", "cultist_acolyte", "imp", "bat").
		AllowEvents(events, "Bandit Camp", "Portal").
		DoTasks(&c.Character)
}
