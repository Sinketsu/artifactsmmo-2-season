package strategy

import (
	"time"

	"github.com/Sinketsu/artifactsmmo/internal/generic"
)

type emptyStrategy struct {
}

func EmptyStrategy() *emptyStrategy {
	return &emptyStrategy{}
}

func (s *emptyStrategy) Do(c *generic.Character) error {
	time.Sleep(1 * time.Second)
	return nil
}
