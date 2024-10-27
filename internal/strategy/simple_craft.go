package strategy

import (
	"fmt"
	"math"
	"time"

	oas "github.com/Sinketsu/artifactsmmo/gen/oas"
	"github.com/Sinketsu/artifactsmmo/internal/bank"
	"github.com/Sinketsu/artifactsmmo/internal/generic"
)

type SimpleCraftStrategy struct {
	craft       []string
	recycle     []string
	deposit     []string
	depositGold bool
	sell        []string
}

// NewTasksFightStrategy returns strategy that will just try spend all resources to craft one item
func NewSimpleCraftStrategy() *SimpleCraftStrategy {
	return &SimpleCraftStrategy{}
}

// Craft sets an item to craft. Required
func (s *SimpleCraftStrategy) Craft(item ...string) *SimpleCraftStrategy {
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

// Deposit sets items to deposit in Bank after craft
func (s *SimpleCraftStrategy) Deposit(items ...string) *SimpleCraftStrategy {
	s.deposit = items
	return s
}

// DepositGold sets allowment to deposit all gold from inventory
func (s *SimpleCraftStrategy) DepositGold() *SimpleCraftStrategy {
	s.depositGold = true
	return s
}

func (s *SimpleCraftStrategy) Do(c *generic.Character) error {
	if s.craft == nil {
		return fmt.Errorf("craft items are not set")
	}

	if len(s.recycle) > 0 {
		if err := c.MacroRecycleAll(s.recycle...); err != nil {
			return fmt.Errorf("recycle: %w", err)
		}
	}

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

	bankItems, err := c.Bank().Items()
	if err != nil {
		return fmt.Errorf("list bank items: %w", err)
	}

	for _, itemCode := range s.craft {
		item, err := c.GetItem(itemCode, true)
		if err != nil {
			return fmt.Errorf("get item: %w", err)
		}

		if !item.Craft.Set {
			c.Logger().Warn("item " + itemCode + " is not craftable!")
			continue
		}

		if s.canCraft(c, item.Craft.Value.CraftSchema.Items, bankItems) {
			return s.craftHelper(c, item)
		}
	}

	time.Sleep(1 * time.Second)
	return nil
}

func (s *SimpleCraftStrategy) craftHelper(c *generic.Character, item oas.SingleItemSchemaItem) error {
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
			} else if withdraw < 0 {
				if err := c.MacroDeposit(ing.Code, -withdraw); err != nil {
					return fmt.Errorf("deposit: %w", err)
				}
			}
		}

		if err := c.MacroCraft(item.Code, trueCount); err != nil {
			return fmt.Errorf("craft: %w", err)
		}

		return nil
	}

	time.Sleep(1 * time.Second)
	return nil
}

func (s *SimpleCraftStrategy) canCraft(c *generic.Character, ingridients []oas.SimpleItemSchema, bank []bank.Item) bool {
	result := math.MaxInt
	for _, req := range ingridients {
		inBank := 0
		for _, bankItem := range bank {
			if bankItem.Code == req.Code {
				inBank = bankItem.Quantity
				break
			}
		}

		q := (inBank + c.InInventory(req.Code)) / req.Quantity
		if q < result {
			result = q
		}
	}

	return result > 0
}
