package strategy

import (
	"fmt"

	"github.com/Sinketsu/artifactsmmo/internal/generic"
)

type SimpleFightStrategy struct {
	fight         string
	sell          []string
	deposit       []string
	depositGold   bool
	allowedEvents []string

	// cache state for current monster
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

// Deposit sets items to deposit in Bank when inventory is full
func (s *SimpleFightStrategy) Deposit(items ...string) *SimpleFightStrategy {
	s.deposit = items
	return s
}

// DepositGold sets allowment to deposit all gold from inventory
func (s *SimpleFightStrategy) DepositGold() *SimpleFightStrategy {
	s.depositGold = true
	return s
}

// AllowEvents sets list of allowed events. When event will be active - fight against event monsters, else fight against usual monster, setted in Fight
func (s *SimpleFightStrategy) AllowEvents(names ...string) *SimpleFightStrategy {
	s.allowedEvents = names
	return s
}

func (s *SimpleFightStrategy) Do(c *generic.Character) error {
	if s.fight == "" {
		return fmt.Errorf("monster not set")
	}

	if c.InventoryItemCount() == c.Data().InventoryMaxItems {
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

	return s.fightHelper(c, s.fight, true)
}

func (s *SimpleFightStrategy) fightHelper(c *generic.Character, code string, cachable bool) error {
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
