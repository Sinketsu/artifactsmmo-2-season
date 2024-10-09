package generic

import (
	"fmt"
	"math"

	api "github.com/Sinketsu/artifactsmmo/gen/oas"
)

func (c *Character) MacroWithdraw(code string, quantity int) error {
	err := c.Move(4, 1) // Bank
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	return c.Withdraw(code, quantity)
}

func (c *Character) MacroSellAll(codes ...string) error {
	err := c.Move(5, 1) // Grand Exchange
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	for _, code := range codes {
		for c.InInventory(code) > 0 {
			ge, err := c.GetGEItem(code)
			if err != nil {
				return fmt.Errorf("get ge item: %w", err)
			}

			sellCount := int(math.Min(float64(c.InInventory(code)), float64(ge.MaxQuantity)))

			_, err = c.Sell(code, sellCount, ge.SellPrice.Value)
			if err != nil {
				return fmt.Errorf("sell: %w", err)
			}
		}
	}

	return nil
}

func (c *Character) MacroDepositAll(codes ...string) error {
	err := c.Move(4, 1) // Bank
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	for _, code := range codes {
		inInventory := c.InInventory(code)
		if inInventory == 0 {
			continue
		}

		if err := c.Deposit(code, inInventory); err != nil {
			return fmt.Errorf("deposit: %w", err)
		}
	}

	return nil
}

func (c *Character) MacroGather(code string) error {
	if c.gatherData.Code != code {
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

		switch resource.Skill {
		case api.ResourceSchemaSkillFishing:
			if c.InInventory("spruce_fishing_rod") > 0 {
				// TODO make helper to choose best tool
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

		c.gatherData.X = tiles[0].X
		c.gatherData.Y = tiles[0].Y
		c.gatherData.Code = code
	}

	err := c.Move(c.gatherData.X, c.gatherData.Y)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	_, err = c.Gather()
	if err != nil {
		return fmt.Errorf("gather: %w", err)
	}

	return nil
}

func (c *Character) MacroFight(monster string) error {
	if c.fightData.Monster != monster {
		tiles, err := c.FindOnMap(monster)
		if err != nil {
			return fmt.Errorf("find on map: %w", err)
		}

		if len(tiles) == 0 {
			return fmt.Errorf("find on map: not found")
		}

		bestGear, err := getBestGearFor(c, monster)
		if err != nil {
			return fmt.Errorf("get best gear: %w", err)
		}

		for _, items := range bestGear {
			if err := c.MacroWear(items); err != nil {
				return fmt.Errorf("wear: %w", err)
			}
		}

		c.fightData.X = tiles[0].X
		c.fightData.Y = tiles[0].Y
		c.fightData.Monster = monster
	}

	err := c.Move(c.fightData.X, c.fightData.Y)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	_, err = c.Fight()
	if err != nil {
		return fmt.Errorf("fight: %w", err)
	}

	return nil
}

func (c *Character) MacroCheckCraftResources(code string) (int, error) {
	item, err := c.GetItem(code)
	if err != nil {
		return 0, fmt.Errorf("get item: %w", err)
	}

	result := math.MaxInt
	for _, req := range item.Craft.Value.CraftSchema.Items {
		q := c.InInventory(req.Code) / req.Quantity
		if q < result {
			result = q
		}
	}

	if result == math.MaxInt {
		return 0, fmt.Errorf("item cannot be crafted")
	}

	return result, nil
}

func (c *Character) MacroCraft(code string, quantity int) error {
	if c.craftData.Code != code {
		item, err := c.GetItem(code)
		if err != nil {
			return fmt.Errorf("get item: %w", err)
		}

		tiles, err := c.FindOnMap(string(item.Craft.Value.CraftSchema.Skill.Value))
		if err != nil {
			return fmt.Errorf("find on map: %w", err)
		}

		if len(tiles) == 0 {
			return fmt.Errorf("find on map: not found")
		}

		c.craftData.X = tiles[0].X
		c.craftData.Y = tiles[0].Y
		c.craftData.Code = code
	}

	err := c.Move(c.craftData.X, c.craftData.Y)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	_, err = c.Craft(code, quantity)
	if err != nil {
		return fmt.Errorf("craft: %w", err)
	}

	return nil
}

func (c *Character) MacroRecycleAll(codes ...string) error {
	for _, code := range codes {
		inInventory := c.InInventory(code)
		if inInventory == 0 {
			continue
		}

		item, err := c.GetItem(code)
		if err != nil {
			return fmt.Errorf("get item: %w", err)
		}

		tiles, err := c.FindOnMap(string(item.Craft.Value.CraftSchema.Skill.Value))
		if err != nil {
			return fmt.Errorf("find on map: %w", err)
		}

		if err := c.Move(tiles[0].X, tiles[0].Y); err != nil {
			return fmt.Errorf("move: %w", err)
		}

		_, err = c.Recycle(code, inInventory)
		if err != nil {
			return fmt.Errorf("recycle: %w", err)
		}
	}

	return nil
}

func (c *Character) MacroWear(items []api.SingleItemSchemaItem) error {
	ringCount := 0

	for _, item := range items {
		switch item.Type {
		case "weapon":
			if c.Data().WeaponSlot == item.Code {
				continue
			}

			if c.Data().WeaponSlot != "" {
				if err := c.UnEquip(c.Data().WeaponSlot, item.Type, 1); err != nil {
					return fmt.Errorf("unequip: %w", err)
				}
			}

			if err := c.Equip(item.Code, item.Type, 1); err != nil {
				return err
			}
		case "helmet":
			if c.Data().HelmetSlot == item.Code {
				continue
			}

			if c.Data().HelmetSlot != "" {
				if err := c.UnEquip(c.Data().HelmetSlot, item.Type, 1); err != nil {
					return fmt.Errorf("unequip: %w", err)
				}
			}

			if err := c.Equip(item.Code, item.Type, 1); err != nil {
				return err
			}
		case "body_armor":
			if c.Data().BodyArmorSlot == item.Code {
				continue
			}

			if c.Data().BodyArmorSlot != "" {
				if err := c.UnEquip(c.Data().BodyArmorSlot, item.Type, 1); err != nil {
					return fmt.Errorf("unequip: %w", err)
				}
			}

			if err := c.Equip(item.Code, item.Type, 1); err != nil {
				return err
			}
		case "leg_armor":
			if c.Data().LegArmorSlot == item.Code {
				continue
			}

			if c.Data().LegArmorSlot != "" {
				if err := c.UnEquip(c.Data().LegArmorSlot, item.Type, 1); err != nil {
					return fmt.Errorf("unequip: %w", err)
				}
			}

			if err := c.Equip(item.Code, item.Type, 1); err != nil {
				return err
			}
		case "shield":
			if c.Data().ShieldSlot == item.Code {
				continue
			}

			if c.Data().ShieldSlot != "" {
				if err := c.UnEquip(c.Data().ShieldSlot, item.Type, 1); err != nil {
					return fmt.Errorf("unequip: %w", err)
				}
			}

			if err := c.Equip(item.Code, item.Type, 1); err != nil {
				return err
			}
		case "boots":
			if c.Data().BootsSlot == item.Code {
				continue
			}

			if c.Data().BootsSlot != "" {
				if err := c.UnEquip(c.Data().BootsSlot, item.Type, 1); err != nil {
					return fmt.Errorf("unequip: %w", err)
				}
			}

			if err := c.Equip(item.Code, item.Type, 1); err != nil {
				return err
			}
		case "amulet":
			if c.Data().AmuletSlot == item.Code {
				continue
			}

			if c.Data().AmuletSlot != "" {
				if err := c.UnEquip(c.Data().AmuletSlot, item.Type, 1); err != nil {
					return fmt.Errorf("unequip: %w", err)
				}
			}

			if err := c.Equip(item.Code, item.Type, 1); err != nil {
				return err
			}
		case "ring":
			if ringCount == 0 {
				ringCount++

				if c.Data().Ring1Slot == item.Code {
					continue
				}

				if c.Data().Ring1Slot != "" {
					if err := c.UnEquip(c.Data().Ring1Slot, "ring1", 1); err != nil {
						return fmt.Errorf("unequip: %w", err)
					}
				}

				if err := c.Equip(item.Code, "ring1", 1); err != nil {
					return fmt.Errorf("equip: %w", err)
				}
			} else {
				if c.Data().Ring2Slot == item.Code {
					continue
				}

				if c.Data().Ring2Slot != "" {
					if err := c.UnEquip(c.Data().Ring2Slot, "ring2", 1); err != nil {
						return fmt.Errorf("unequip: %w", err)
					}
				}

				if err := c.Equip(item.Code, "ring2", 1); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("unknown type: %s", item.Type)
		}
	}
	return nil
}

func (c *Character) MacroCompleteMonsterTask() error {
	err := c.Move(1, 2) // Task master monsters
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	reward, err := c.CompleteTask()
	if err != nil {
		return fmt.Errorf("complete task: %w", err)
	}

	c.Log("completed task, got", reward.Quantity, reward.Code)
	return nil
}

func (c *Character) MacroNewMonsterTask() error {
	err := c.Move(1, 2) // Task master monsters
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	task, err := c.AcceptNewTask()
	if err != nil {
		return fmt.Errorf("complete task: %w", err)
	}

	c.Log("accept new task:", task.Total, task.Code)
	return nil
}

func (c *Character) MacroCancelMonsterTask() error {
	err := c.Move(1, 2) // Task master monsters
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	err = c.CancelTask()
	if err != nil {
		return fmt.Errorf("cancel task: %w", err)
	}

	c.Log("cancelled task")
	return nil
}
