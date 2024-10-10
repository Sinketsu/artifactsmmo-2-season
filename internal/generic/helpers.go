package generic

import (
	"context"
	"fmt"
	"strings"

	api "github.com/Sinketsu/artifactsmmo/gen/oas"
	combinations "github.com/mxschmitt/golang-combinations"
)

func gearCombinations(idx int, all ...[][]api.SingleItemSchemaItem) [][]api.SingleItemSchemaItem {
	result := make([][]api.SingleItemSchemaItem, 0)

	if idx == len(all)-1 {
		return all[idx]
	}

	if len(all[idx]) == 0 {
		return gearCombinations(idx+1, all...)
	}

	for _, slotCombinations := range all[idx] {
		remainingSet := gearCombinations(idx+1, all...)
		if len(remainingSet) == 0 {
			result = append(result, slotCombinations)
			continue
		}

		for _, other := range remainingSet {
			newSet := []api.SingleItemSchemaItem{}
			newSet = append(newSet, slotCombinations...)
			newSet = append(newSet, other...)

			result = append(result, newSet)
		}
	}

	return result
}

func (c *Character) GetBestGearFor(monsterCode string) ([]api.SingleItemSchemaItem, error) {
	monster, err := c.GetMonster(monsterCode, true)
	if err != nil {
		return nil, fmt.Errorf("get monster: %w", err)
	}

	var weapons, helmets, shields, bodyArmors, legsArmors, boots, rings, amulets []api.SingleItemSchemaItem
	for _, i := range c.Data().Inventory {
		if i.Code == "" {
			continue
		}

		item, err := c.GetItem(i.Code, true)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", i.Code, err)
		}

		switch item.Type {
		case "weapon":
			weapons = append(weapons, item)
		case "helmet":
			helmets = append(helmets, item)
		case "shield":
			shields = append(shields, item)
		case "body_armor":
			bodyArmors = append(bodyArmors, item)
		case "leg_armor":
			legsArmors = append(legsArmors, item)
		case "boots":
			boots = append(boots, item)
		case "ring":
			rings = append(rings, item)
		case "amulet":
			amulets = append(amulets, item)
		}
	}

	if c.Data().WeaponSlot != "" {
		item, err := c.GetItem(c.Data().WeaponSlot, true)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().WeaponSlot, err)
		}
		weapons = append(weapons, item)
	}
	if c.Data().HelmetSlot != "" {
		item, err := c.GetItem(c.Data().HelmetSlot, true)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().HelmetSlot, err)
		}
		helmets = append(helmets, item)
	}
	if c.Data().BodyArmorSlot != "" {
		item, err := c.GetItem(c.Data().BodyArmorSlot, true)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().BodyArmorSlot, err)
		}
		bodyArmors = append(bodyArmors, item)
	}
	if c.Data().LegArmorSlot != "" {
		item, err := c.GetItem(c.Data().LegArmorSlot, true)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().LegArmorSlot, err)
		}
		legsArmors = append(legsArmors, item)
	}
	if c.Data().BootsSlot != "" {
		item, err := c.GetItem(c.Data().BootsSlot, true)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().BootsSlot, err)
		}
		boots = append(boots, item)
	}
	if c.Data().ShieldSlot != "" {
		item, err := c.GetItem(c.Data().ShieldSlot, true)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().ShieldSlot, err)
		}
		shields = append(shields, item)
	}
	if c.Data().Ring1Slot != "" {
		item, err := c.GetItem(c.Data().Ring1Slot, true)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().Ring1Slot, err)
		}
		rings = append(rings, item)
	}
	if c.Data().Ring2Slot != "" {
		item, err := c.GetItem(c.Data().Ring2Slot, true)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().Ring2Slot, err)
		}
		rings = append(rings, item)
	}
	if c.Data().AmuletSlot != "" {
		item, err := c.GetItem(c.Data().AmuletSlot, true)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().AmuletSlot, err)
		}
		amulets = append(amulets, item)
	}

	weaponCombinations := combinations.Combinations(weapons, 1)
	bodyArmorCombinations := combinations.Combinations(bodyArmors, 1)
	legsArmorCombinations := combinations.Combinations(legsArmors, 1)
	shieldCombinations := combinations.Combinations(shields, 1)
	amuletCombinations := combinations.Combinations(amulets, 1)
	bootsCombinations := combinations.Combinations(boots, 1)
	ringsCombinations := combinations.Combinations(rings, 2)
	helmetCombinations := combinations.Combinations(helmets, 1)

	bestScore := 0.0
	bestGear := []api.SingleItemSchemaItem{}

	setCombinations := gearCombinations(0, weaponCombinations, bodyArmorCombinations, legsArmorCombinations,
		shieldCombinations, amuletCombinations, bootsCombinations, ringsCombinations, helmetCombinations)

	for _, combination := range setCombinations {
		totalEffects := make(map[string]int)

		for _, item := range combination {
			for _, e := range item.Effects {
				totalEffects[e.Name] += e.Value
			}
		}

		score := float64(totalEffects["attack_earth"])*(1-float64(monster.ResEarth)/100)*(1+float64(totalEffects["dmg_earth"])/100) +
			float64(totalEffects["attack_water"])*(1-float64(monster.ResWater)/100)*(1+float64(totalEffects["dmg_water"])/100) +
			float64(totalEffects["attack_fire"])*(1-float64(monster.ResFire)/100)*(1+float64(totalEffects["dmg_fire"])/100) +
			float64(totalEffects["attack_air"])*(1-float64(monster.ResAir)/100)*(1+float64(totalEffects["dmg_air"])/100)
		if score > bestScore {
			bestScore = score
			bestGear = combination
		}
	}

	bestGearCodes := make([]string, 0, len(bestGear))
	for _, gear := range bestGear {
		bestGearCodes = append(bestGearCodes, gear.Code)
	}
	c.Log("found best gear for monster", monster.Code, "with effective dmg:", bestScore, "[", strings.Join(bestGearCodes, ", "), "]")

	return bestGear, nil
}

func (c *Character) InventoryItemCount() int {
	count := 0
	for _, item := range c.data.Inventory {
		count += item.Quantity
	}

	return count
}

func (c *Character) EmptyInventorySlots() int {
	count := 0
	for _, item := range c.data.Inventory {
		if item.Code == "" {
			count++
		}
	}

	return count
}

func (c *Character) InInventory(code string) int {
	for _, item := range c.data.Inventory {
		if item.Code == code {
			return item.Quantity
		}
	}

	return 0
}

func (c *Character) InBank(code string) (int, error) {
	res, err := c.cli.GetBankItemsMyBankItemsGet(context.Background(), api.GetBankItemsMyBankItemsGetParams{ItemCode: api.NewOptString(code)})
	if err != nil {
		return 0, err
	}

	if len(res.Data) == 0 {
		return 0, nil
	}

	return res.Data[0].Quantity, nil
}

func (c *Character) InventoryIsFull() bool {
	return c.InventoryItemCount() == c.Data().InventoryMaxItems || c.EmptyInventorySlots() == 0
}
