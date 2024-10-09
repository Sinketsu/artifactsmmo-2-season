package strategy

import "github.com/Sinketsu/artifactsmmo/internal/generic"

type Strategy interface {
	Do(c *generic.Character) error
}

type gatherInfo struct {
	Code string
	X    int
	Y    int
}

type fightInfo struct {
	Code string
	X    int
	Y    int
}
