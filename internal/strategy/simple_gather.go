package strategy

import (
	"fmt"

	"github.com/Sinketsu/artifactsmmo/internal/generic"
)

type SimpleGatherStrategy struct {
	gather        string
	sell          []string
	deposit       []string
	depositGold   bool
	craft         string
	allowedEvents []string

	// cache state for current resource
	currentResource string
}

// NewTasksFightStrategy returns strategy that will just simple gather one resource
func NewSimpleGatherStrategy() *SimpleGatherStrategy {
	return &SimpleGatherStrategy{}
}

// Gather sets a resource to gather. Required
func (s *SimpleGatherStrategy) Gather(resourceSpot string) *SimpleGatherStrategy {
	s.gather = resourceSpot
	return s
}

// Sell sets items to sell in GE when inventory is full
func (s *SimpleGatherStrategy) Sell(items ...string) *SimpleGatherStrategy {
	s.sell = items
	return s
}

// Deposit sets items to deposit in Bank when inventory is full
func (s *SimpleGatherStrategy) Deposit(items ...string) *SimpleGatherStrategy {
	s.deposit = items
	return s
}

// DepositGold sets allowment to deposit all gold from inventory
func (s *SimpleGatherStrategy) DepositGold() *SimpleGatherStrategy {
	s.depositGold = true
	return s
}

// Craft sets items to try craft when inventory is full. If no such resources to craft - it will be ignored. Runs before sell or bank triggers
func (s *SimpleGatherStrategy) Craft(item string) *SimpleGatherStrategy {
	s.craft = item
	return s
}

// AllowEvents sets list of allowed events. When event will be active - gather event resources, else gather usual resource,  setted in Gather
func (s *SimpleGatherStrategy) AllowEvents(names ...string) *SimpleGatherStrategy {
	s.allowedEvents = names
	return s
}

func (s *SimpleGatherStrategy) Do(c *generic.Character) error {
	if s.gather == "" {
		return fmt.Errorf("gather resource not set")
	}

	if c.InventoryIsFull() {
		c.Logger().Info("inventory is full - going to craft, bank, GE...")

		if s.craft != "" {
			q, err := c.MacroCheckCraftResources(s.craft)
			if err != nil {
				return fmt.Errorf("check craft resources: %w", err)
			}

			if q > 0 {
				if err := c.MacroCraft(s.craft, q); err != nil {
					return fmt.Errorf("craft: %w", err)
				}
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
	}

	for _, eventName := range s.allowedEvents {
		if event := c.Events().Get(eventName); event != nil {
			return s.gatherHelper(c, event.Map.Content.MapContentSchema.Code, false)
		}
	}

	return s.gatherHelper(c, s.gather, true)
}

func (s *SimpleGatherStrategy) gatherHelper(c *generic.Character, code string, cachable bool) error {
	if s.currentResource != code {
		c.Bank().Lock()

		gear, err := c.GetBestGatherGearFor(code)
		if err != nil {
			c.Bank().Unlock()
			return fmt.Errorf("get best gear: %w", err)
		}

		err = c.MacroWear(gear)
		if err != nil {
			c.Bank().Unlock()
			return fmt.Errorf("wear: %w", err)
		}

		c.Bank().Unlock()
		s.currentResource = code
	}

	return c.MacroGather(code, cachable)
}
