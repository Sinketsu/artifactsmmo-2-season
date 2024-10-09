package strategy

import (
	"fmt"
	"math"
	"time"

	"github.com/Sinketsu/artifactsmmo/internal/generic"
)

type SimpleCraftStrategy struct {
	craft   string
	recycle []string
	bank    []string
	sell    []string
}

// NewTasksFightStrategy returns strategy that will just try spend all resources to craft one item
func NewSimpleCraftStrategy() *SimpleCraftStrategy {
	return &SimpleCraftStrategy{}
}

// Craft sets an item to craft. Required
func (s *SimpleCraftStrategy) Craft(item string) *SimpleCraftStrategy {
	s.craft = item
	return s
}

// Recycle sets items to recycle after craft
func (s *SimpleCraftStrategy) Recycle(items ...string) *SimpleCraftStrategy {
	s.recycle = items
	return s
}

// Sell sets items to sell in GE after craft
func (s *SimpleCraftStrategy) Sell(items ...string) *SimpleCraftStrategy {
	s.sell = items
	return s
}

// Bank sets items to deposit in Bank after craft
func (s *SimpleCraftStrategy) Bank(items ...string) *SimpleCraftStrategy {
	s.bank = items
	return s
}

func (s *SimpleCraftStrategy) Do(c *generic.Character) error {
	if s.craft == "" {
		return fmt.Errorf("craft item not set")
	}

	if len(s.recycle) > 0 {
		if err := c.MacroRecycleAll(s.recycle...); err != nil {
			return fmt.Errorf("recycle: %w", err)
		}
	}

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

	return s.crafttHelper(c)
}

func (s *SimpleCraftStrategy) crafttHelper(c *generic.Character) error {
	item, err := c.GetItem(s.craft, true)
	if err != nil {
		return fmt.Errorf("get item: %w", err)
	}

	if !item.Craft.Set {
		return fmt.Errorf("item is not craftable")
	}

	ingridients := item.Craft.Value.CraftSchema.Items

	// main craft logic
	space := c.Data().InventoryMaxItems - c.InventoryItemCount()
	ingridientCount := 0
	minAvailableCount := 99999999
	for _, ing := range ingridients {
		inventory := c.InInventory(ing.Code)
		bank, _ := c.InBank(ing.Code)

		count := (inventory + bank) / ing.Quantity
		if count < minAvailableCount {
			minAvailableCount = count
		}

		space += inventory
		ingridientCount += ing.Quantity
	}

	if minAvailableCount == 0 {
		time.Sleep(1 * time.Second)
		return nil
	}

	spaceAvailableCount := space / ingridientCount
	trueCount := int(math.Min(float64(minAvailableCount), float64(spaceAvailableCount)))

	if trueCount > 0 {
		for _, ing := range ingridients {
			inventory := c.InInventory(ing.Code)

			withdraw := trueCount*ing.Quantity - inventory

			if withdraw > 0 {
				if err := c.MacroWithdraw(ing.Code, withdraw); err != nil {
					return fmt.Errorf("withdraw: %w", err)
				}
			}
		}

		if err := c.MacroCraft(s.craft, trueCount); err != nil {
			return fmt.Errorf("craft: %w", err)
		}
	}

	time.Sleep(1 * time.Second)
	return nil
}
