package strategy

import (
	"fmt"
	"math"
	"time"

	api "github.com/Sinketsu/artifactsmmo/gen/oas"
	"github.com/Sinketsu/artifactsmmo/internal/generic"
)

type craftData struct {
	code        string
	ingridients []api.SimpleItemSchema
}

type SimpleCraftStrategy struct {
	craft   string
	recycle []string
	bank    []string
	sell    []string

	craftData craftData
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

	if s.craftData.code != s.craft {
		item, err := c.GetItem(s.craft)
		if err != nil {
			return fmt.Errorf("get item: %w", err)
		}

		if !item.Craft.Set {
			return fmt.Errorf("item is not craftable")
		}

		s.craftData.ingridients = item.Craft.Value.CraftSchema.Items
		s.craftData.code = s.craft
	}

	space := c.Data().InventoryMaxItems - c.InventoryItemCount()
	ingridientCount := 0
	minAvailableCount := 99999999
	for _, ing := range s.craftData.ingridients {
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
		for _, ing := range s.craftData.ingridients {
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
