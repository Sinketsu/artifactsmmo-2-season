package ereshkigal

import (
	"context"
	"time"

	"github.com/Sinketsu/artifactsmmo/internal/generic"
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
	return generic.NewSimpleFightStrategy().
		Fight("yellow_slime").
		Bank("yellow_slimeball", "green_slimeball", "feather").
		Do(&c.Character)
}
