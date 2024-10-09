package strategy

import (
	"fmt"
	"slices"

	"github.com/Sinketsu/artifactsmmo/internal/events"
	"github.com/Sinketsu/artifactsmmo/internal/generic"
)

type tasksFightStrategy struct {
	sell          []string
	bank          []string
	cancelTasks   []string
	events        *events.Service
	allowedEvents []string

	// cache state for current mosnter
	currentMonster string
}

// NewTasksFightStrategy returns strategy that will do monster task quests
func NewTasksFightStrategy() *tasksFightStrategy {
	return &tasksFightStrategy{}
}

// Sell sets items to sell in GE when inventory is full
func (s *tasksFightStrategy) Sell(items ...string) *tasksFightStrategy {
	s.sell = items
	return s
}

// Bank sets items to deposit in Bank when inventory is full
func (s *tasksFightStrategy) Bank(items ...string) *tasksFightStrategy {
	s.bank = items
	return s
}

// CancelTasks sets a list of blacklisted tasks - which we do not want to do (too weak or ineffective)
func (s *tasksFightStrategy) CancelTasks(monsters ...string) *tasksFightStrategy {
	s.cancelTasks = monsters
	return s
}

// AllowEvents sets list of allowed events. When event will be active - fight against event monsters, else fight against usual monster, setted in Fight
func (s *tasksFightStrategy) AllowEvents(events *events.Service, names ...string) *tasksFightStrategy {
	s.allowedEvents = names
	s.events = events
	return s
}

func (s *tasksFightStrategy) Do(c *generic.Character) error {
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
	}

	if s.events != nil {
		for _, eventName := range s.allowedEvents {
			if event := s.events.Get(eventName); event != nil {
				return s.fightHelper(c, event.Map.Content.MapContentSchema.Code, false)
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

	return s.fightHelper(c, c.Data().Task, true)
}

func (s *tasksFightStrategy) fightHelper(c *generic.Character, code string, cachable bool) error {
	tile, err := c.FindOnMap(code, cachable)
	if err != nil {
		return fmt.Errorf("find on map: %w", err)
	}

	if s.currentMonster != code {
		bestGear, err := c.GetBestGearFor(code)
		if err != nil {
			return fmt.Errorf("get best gear: %w", err)
		}

		for _, items := range bestGear {
			if err := c.MacroWear(items); err != nil {
				return fmt.Errorf("wear: %w", err)
			}
		}

		s.currentMonster = code
	}

	return c.MacroFight(tile.X, tile.Y)
}
