package strategy

import (
	"fmt"

	"github.com/Sinketsu/artifactsmmo/internal/events"
	"github.com/Sinketsu/artifactsmmo/internal/generic"
)

type SimpleFightStrategy struct {
	fight         string
	sell          []string
	bank          []string
	events        *events.Service
	allowedEvents []string

	// cache state for current mosnter
	currentMonster string
}

// NewTasksFightStrategy returns strategy that will just simple fight against one monster
func NewSimpleFightStrategy() *SimpleFightStrategy {
	return &SimpleFightStrategy{}
}

// Fight sets a monster to fight. Required
func (s *SimpleFightStrategy) Fight(monster string) *SimpleFightStrategy {
	s.fight = monster
	return s
}

// Sell sets items to sell in GE when inventory is full
func (s *SimpleFightStrategy) Sell(items ...string) *SimpleFightStrategy {
	s.sell = items
	return s
}

// Bank sets items to deposit in Bank when inventory is full
func (s *SimpleFightStrategy) Bank(items ...string) *SimpleFightStrategy {
	s.bank = items
	return s
}

// AllowEvents sets list of allowed events. When event will be active - fight against event monsters, else fight against usual monster, setted in Fight
func (s *SimpleFightStrategy) AllowEvents(events *events.Service, names ...string) *SimpleFightStrategy {
	s.allowedEvents = names
	s.events = events
	return s
}

func (s *SimpleFightStrategy) Do(c *generic.Character) error {
	if s.fight == "" {
		return fmt.Errorf("monster not set")
	}

	if c.InventoryItemCount() == c.Data().InventoryMaxItems {
		c.Log("inventory is full - going to bank, GE...")

		if len(s.bank) > 0 {
			if err := c.MacroDepositAll(s.bank...); err != nil {
				return fmt.Errorf("deposit: %w", err)
			}
		}

		if len(s.sell) > 0 {
			if err := c.MacroSellAll(s.sell...); err != nil {
				return fmt.Errorf("sell: %w", err)
			}
		}
	}

	if s.events != nil {
		for _, eventName := range s.allowedEvents {
			if event := s.events.Get(eventName); event != nil {
				return s.fightHelper(c, event.Map.Content.MapContentSchema.Code, false)
			}
		}
	}

	return s.fightHelper(c, s.fight, true)
}

func (s *SimpleFightStrategy) fightHelper(c *generic.Character, code string, cachable bool) error {
	tile, err := c.FindOnMap(code, cachable)
	if err != nil {
		return fmt.Errorf("find on map: %w", err)
	}

	if s.currentMonster != code {
		bestGear, err := c.GetBestGearFor(code)
		if err != nil {
			return fmt.Errorf("get best gear: %w", err)
		}

		if err := c.MacroWear(bestGear); err != nil {
			return fmt.Errorf("wear: %w", err)
		}

		s.currentMonster = code
	}

	return c.MacroFight(tile.X, tile.Y)
}
