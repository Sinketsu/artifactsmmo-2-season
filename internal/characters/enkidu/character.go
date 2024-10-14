package enkidu

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
			err := c.do()
			if err != nil {
				c.Log(err)
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (c *Character) do() error {

	c.setStrategy("waiting for resources...", strategy.EmptyStrategy())

	return c.strategy.Do(&c.Character)
}

func (c *Character) setStrategy(what string, newStrategy strategy.Strategy) {
	if c.what != what {
		c.Log("change strategy:", what)
		c.strategy = newStrategy
		c.what = what
	}
}
