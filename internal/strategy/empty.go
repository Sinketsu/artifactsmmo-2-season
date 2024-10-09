package strategy

import (
	"github.com/Sinketsu/artifactsmmo/internal/generic"
)

type emptyStrategy struct {
}

func EmptyStrategy() *emptyStrategy {
	return &emptyStrategy{}
}

func (s *emptyStrategy) Do(c *generic.Character) error {
	return nil
}

func (s *emptyStrategy) DoTasks(c *generic.Character) error {
	return nil
}
