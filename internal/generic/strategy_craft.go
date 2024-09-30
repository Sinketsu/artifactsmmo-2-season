package generic

import (
	"fmt"
	"time"
)

type SimpleCraftStrategy struct {
	craft   string
	recycle []string
	bank    []string
	sell    []string
}

func NewSimpleCraftStrategy() *SimpleCraftStrategy {
	return &SimpleCraftStrategy{}
}

func (s *SimpleCraftStrategy) Craft(item string) *SimpleCraftStrategy {
	s.craft = item
	return s
}

func (s *SimpleCraftStrategy) Recycle(items ...string) *SimpleCraftStrategy {
	s.recycle = items
	return s
}

func (s *SimpleCraftStrategy) Sell(items ...string) *SimpleCraftStrategy {
	s.sell = items
	return s
}

func (s *SimpleCraftStrategy) Bank(items ...string) *SimpleCraftStrategy {
	s.bank = items
	return s
}

func (s *SimpleCraftStrategy) Do(c *Character) error {
	if s.craft == "" {
		return fmt.Errorf("craft item not set")
	}

	q, err := c.MacroCheckCraftResources(s.craft)
	if err != nil {
		return fmt.Errorf("check craft resources: %w", err)
	}

	if q > 0 {
		if err := c.MacroCraft(s.craft, q); err != nil {
			return fmt.Errorf("craft: %w", err)
		}
		return nil
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

		if len(s.recycle) > 0 {
			if err := c.MacroRecycleAll(s.recycle...); err != nil {
				return fmt.Errorf("recycle: %w", err)
			}
		}

		return nil
	}

	item, err := c.GetItem(s.craft)
	if err != nil {
		return fmt.Errorf("get item: %w", err)
	}

	tryWithdraw := false

	for _, craftItem := range item.Craft.Value.CraftSchema.Items {
		if c.InInventory(craftItem.Code) >= craftItem.Quantity {
			continue
		}

		tryWithdraw = true
		if err := c.Move(4, 1); err != nil { // Bank
			return fmt.Errorf("move: %w", err)
		}

		need := craftItem.Quantity - c.InInventory(craftItem.Code)
		inBank, _ := c.InBank(craftItem.Code)
		if inBank >= need {
			if err := c.Withdraw(craftItem.Code, need); err != nil {
				return fmt.Errorf("withdraw: %w", err)
			}
		}
	}

	if !tryWithdraw {
		time.Sleep(1 * time.Second)
		c.Log("idle...")
	}

	return nil
}
