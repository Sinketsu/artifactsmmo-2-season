package strategy

import (
	"fmt"
	"slices"

	"github.com/Sinketsu/artifactsmmo/internal/events"
	"github.com/Sinketsu/artifactsmmo/internal/generic"
)

type SimpleFightStrategy struct {
	fight string
	sell  []string
	bank  []string

	cancelTasks []string

	events        *events.Service
	allowedEvents []string
}

func NewSimpleFightStrategy() *SimpleFightStrategy {
	return &SimpleFightStrategy{}
}

func (s *SimpleFightStrategy) Fight(monster string) *SimpleFightStrategy {
	s.fight = monster
	return s
}

func (s *SimpleFightStrategy) Sell(items ...string) *SimpleFightStrategy {
	s.sell = items
	return s
}

func (s *SimpleFightStrategy) Bank(items ...string) *SimpleFightStrategy {
	s.bank = items
	return s
}

func (s *SimpleFightStrategy) CancelTasks(monsters ...string) *SimpleFightStrategy {
	s.cancelTasks = monsters
	return s
}

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

		return nil
	}

	if s.events != nil {
		for _, eventName := range s.allowedEvents {
			if event := s.events.Get(eventName); event != nil {
				return c.MacroFight(event.Map.Content.MapContentSchema.Code)
			}
		}
	}

	return c.MacroFight(s.fight)
}

func (s *SimpleFightStrategy) DoTasks(c *generic.Character) error {
	if c.InventoryItemCount() == c.Data().InventoryMaxItems || c.EmptyInventorySlots() == 0 {
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

		return nil
	}

	if s.events != nil {
		for _, eventName := range s.allowedEvents {
			if event := s.events.Get(eventName); event != nil {
				return c.MacroFight(event.Map.Content.MapContentSchema.Code)
			}
		}
	}

	if c.Data().Task == "" {
		if err := c.MacroNewMonsterTask(); err != nil {
			return fmt.Errorf("accept new task: %w", err)
		}
	}

	if c.Data().TaskProgress == c.Data().TaskTotal {
		if err := c.MacroCompleteMonsterTask(); err != nil {
			return fmt.Errorf("complete task: %w", err)
		}
	}

	if slices.Contains(s.cancelTasks, c.Data().Task) {
		if err := c.MacroCancelMonsterTask(); err != nil {
			return fmt.Errorf("cancel task: %w", err)
		}

		return nil
	}

	return c.MacroFight(c.Data().Task)
}
