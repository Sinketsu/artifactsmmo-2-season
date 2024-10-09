package strategy

import (
	"fmt"

	"github.com/Sinketsu/artifactsmmo/internal/events"
	"github.com/Sinketsu/artifactsmmo/internal/generic"
)

type SimpleGatherStrategy struct {
	gather string
	sell   []string
	bank   []string
	craft  string

	events        *events.Service
	allowedEvents []string
}

func NewSimpleGatherStrategy() *SimpleGatherStrategy {
	return &SimpleGatherStrategy{}
}

func (s *SimpleGatherStrategy) Gather(resourceSpot string) *SimpleGatherStrategy {
	s.gather = resourceSpot
	return s
}

func (s *SimpleGatherStrategy) Sell(items ...string) *SimpleGatherStrategy {
	s.sell = items
	return s
}

func (s *SimpleGatherStrategy) Bank(items ...string) *SimpleGatherStrategy {
	s.bank = items
	return s
}

func (s *SimpleGatherStrategy) Craft(item string) *SimpleGatherStrategy {
	s.craft = item
	return s
}

func (s *SimpleGatherStrategy) AllowEvents(events *events.Service, names ...string) *SimpleGatherStrategy {
	s.allowedEvents = names
	s.events = events
	return s
}

func (s *SimpleGatherStrategy) Do(c *generic.Character) error {
	if s.gather == "" {
		return fmt.Errorf("gather resource not set")
	}

	if c.InventoryItemCount() == c.Data().InventoryMaxItems {
		c.Log("inventory is full - going to craft, bank, GE...")

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
				return c.MacroGather(event.Map.Content.MapContentSchema.Code)
			}
		}
	}

	return c.MacroGather(s.gather)
}

func (s *SimpleGatherStrategy) DoTasks(c *generic.Character) error {
	return fmt.Errorf("not supported yet")
}
