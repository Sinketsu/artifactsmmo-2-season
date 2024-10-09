package strategy

import "github.com/Sinketsu/artifactsmmo/internal/generic"

type Strategy interface {
	Do(c *generic.Character) error
}
