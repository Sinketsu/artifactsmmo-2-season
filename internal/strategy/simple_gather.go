package strategy

import (
	"fmt"

	api "github.com/Sinketsu/artifactsmmo/gen/oas"
	"github.com/Sinketsu/artifactsmmo/internal/events"
	"github.com/Sinketsu/artifactsmmo/internal/generic"
)

type SimpleGatherStrategy struct {
	gather        string
	sell          []string
	bank          []string
	craft         string
	events        *events.Service
	allowedEvents []string

	info gatherInfo
}

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

// Bank sets items to deposit in Bank when inventory is full
func (s *SimpleGatherStrategy) Bank(items ...string) *SimpleGatherStrategy {
	s.bank = items
	return s
}

// Craft sets items to try craft when inventory is full. If no such resources to craft - it will be ignored. Runs before sell or bank triggers
func (s *SimpleGatherStrategy) Craft(item string) *SimpleGatherStrategy {
	s.craft = item
	return s
}

// AllowEvents sets list of allowed events. When event will be active - gather event resources, else gather usual resource,  setted in Gather
func (s *SimpleGatherStrategy) AllowEvents(events *events.Service, names ...string) *SimpleGatherStrategy {
	s.allowedEvents = names
	s.events = events
	return s
}

func (s *SimpleGatherStrategy) Do(c *generic.Character) error {
	if s.gather == "" {
		return fmt.Errorf("gather resource not set")
	}

	if c.InventoryIsFull() {
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
	}

	if s.events != nil {
		for _, eventName := range s.allowedEvents {
			if event := s.events.Get(eventName); event != nil {
				return s.gatherHelper(c, event.Map.Content.MapContentSchema.Code)
			}
		}
	}

	return s.gatherHelper(c, s.gather)
}

func (s *SimpleGatherStrategy) gatherHelper(c *generic.Character, code string) error {
	if s.info.Code != code {
		tiles, err := c.FindOnMap(code)
		if err != nil {
			return fmt.Errorf("find on map: %w", err)
		}

		if len(tiles) == 0 {
			return fmt.Errorf("find on map: not found")
		}

		resource, err := c.GetResource(code)
		if err != nil {
			return fmt.Errorf("get resource: %w", err)
		}

		// TODO make helper to choose best tool
		switch resource.Skill {
		case api.ResourceSchemaSkillFishing:
			if c.InInventory("spruce_fishing_rod") > 0 {
				if err := c.MacroWear([]api.SingleItemSchemaItem{{Code: "spruce_fishing_rod", Type: "weapon"}}); err != nil {
					return fmt.Errorf("wear: %w", err)
				}
			}
		case api.ResourceSchemaSkillMining:
			if c.InInventory("iron_pickaxe") > 0 {
				if err := c.MacroWear([]api.SingleItemSchemaItem{{Code: "iron_pickaxe", Type: "weapon"}}); err != nil {
					return fmt.Errorf("wear: %w", err)
				}
			}
		case api.ResourceSchemaSkillWoodcutting:
			if c.InInventory("iron_axe") > 0 {
				if err := c.MacroWear([]api.SingleItemSchemaItem{{Code: "iron_axe", Type: "weapon"}}); err != nil {
					return fmt.Errorf("wear: %w", err)
				}
			}
		}

		s.info.X = tiles[0].X
		s.info.Y = tiles[0].Y
		s.info.Code = code
	}

	return c.MacroGather(s.info.X, s.info.Y)
}
