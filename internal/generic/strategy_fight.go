package generic

import "fmt"

type SimpleFightStrategy struct {
	fight string
	sell  []string
	bank  []string
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

func (s *SimpleFightStrategy) Do(c *Character) error {
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

	return c.MacroFight(s.fight)
}
