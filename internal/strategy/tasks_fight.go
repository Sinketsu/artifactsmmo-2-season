package strategy

import (
	"fmt"
	"slices"

	"github.com/Sinketsu/artifactsmmo/internal/generic"
)

type tasksFightStrategy struct {
	sell          []string
	deposit       []string
	depositGold   bool
	cancelTasks   []string
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

// Deposit sets items to deposit in Bank when inventory is full
func (s *tasksFightStrategy) Deposit(items ...string) *tasksFightStrategy {
	s.deposit = items
	return s
}

// DepositGold sets allowment to deposit all gold from inventory
func (s *tasksFightStrategy) DepositGold() *tasksFightStrategy {
	s.depositGold = true
	return s
}

// CancelTasks sets a list of blacklisted tasks - which we do not want to do (too weak or ineffective)
func (s *tasksFightStrategy) CancelTasks(monsters ...string) *tasksFightStrategy {
	s.cancelTasks = monsters
	return s
}

// AllowEvents sets list of allowed events. When event will be active - fight against event monsters, else fight against usual monster, setted in Fight
func (s *tasksFightStrategy) AllowEvents(names ...string) *tasksFightStrategy {
	s.allowedEvents = names
	return s
}

func (s *tasksFightStrategy) Do(c *generic.Character) error {
	// we need some space for complete task
	if c.Data().InventoryMaxItems-c.InventoryItemCount() < 10 || c.EmptyInventorySlots() == 0 {
		c.Logger().Info("inventory is full - going to bank, GE...")

		if len(s.deposit) > 0 {
			if err := c.MacroDepositAll(s.deposit...); err != nil {
				return fmt.Errorf("deposit: %w", err)
			}
		}

		if s.depositGold {
			if err := c.MacroDepositGold(c.Data().Gold); err != nil {
				return fmt.Errorf("deposit gold: %w", err)
			}
		}

		if len(s.sell) > 0 {
			if err := c.MacroSellAll(s.sell...); err != nil {
				return fmt.Errorf("sell: %w", err)
			}
		}
	}

	for _, eventName := range s.allowedEvents {
		if event := c.Events().Get(eventName); event != nil {
			return s.fightHelper(c, event.Map.Content.MapContentSchema.Code, false)
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
	if s.currentMonster != code {
		c.Bank().Lock()

		bestGear, err := c.GetBestGearFor(code)
		if err != nil {
			c.Bank().Unlock()
			return fmt.Errorf("get best gear: %w", err)
		}

		if err := c.MacroWear(bestGear); err != nil {
			c.Bank().Unlock()
			return fmt.Errorf("wear: %w", err)
		}

		c.Bank().Unlock()
		s.currentMonster = code
	}

	return c.MacroFight(code, cachable)
}
