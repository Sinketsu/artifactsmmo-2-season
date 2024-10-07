package enkidu

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
	// return generic.NewSimpleFightStrategy().
	// 	Fight("chicken").
	// 	Bank("golden_egg", "feather").
	// 	Sell("raw_chicken", "egg").
	// 	Do(&c.Character)
	return strategy.NewSimpleCraftStrategy().
		Craft("cooked_trout").
		Sell("cooked_trout").
		Do(&c.Character)
	// return generic.NewSimpleCraftStrategy().
	// 	Craft("iron_boots").
	// 	Recycle("iron_boots").
	// 	Do(&c.Character)
	// return nil
}
