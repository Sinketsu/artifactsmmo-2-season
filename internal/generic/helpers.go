package generic

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	oas "github.com/Sinketsu/artifactsmmo/gen/oas"
	combinations "github.com/mxschmitt/golang-combinations"
)

func gearCombinations(f func(items []oas.SingleItemSchemaItem), weaponCombinations, bodyArmorCombinations, legsArmorCombinations,
	shieldCombinations, amuletCombinations, bootsCombinations, ringsCombinations, helmetCombinations [][]oas.SingleItemSchemaItem) {

	var items = make([]oas.SingleItemSchemaItem, 9)
	var empty = [][]oas.SingleItemSchemaItem{{{}}}

	if len(weaponCombinations) == 0 {
		weaponCombinations = empty
	}
	if len(bodyArmorCombinations) == 0 {
		bodyArmorCombinations = empty
	}
	if len(legsArmorCombinations) == 0 {
		legsArmorCombinations = empty
	}
	if len(shieldCombinations) == 0 {
		shieldCombinations = empty
	}
	if len(amuletCombinations) == 0 {
		amuletCombinations = empty
	}
	if len(bootsCombinations) == 0 {
		bootsCombinations = empty
	}
	if len(ringsCombinations) == 0 {
		ringsCombinations = empty
	}
	if len(helmetCombinations) == 0 {
		helmetCombinations = empty
	}

	for _, weapon := range weaponCombinations {
		for _, bodyArmor := range bodyArmorCombinations {
			for _, legArmor := range legsArmorCombinations {
				for _, shield := range shieldCombinations {
					for _, amulet := range amuletCombinations {
						for _, boots := range bootsCombinations {
							for _, ringSet := range ringsCombinations {
								for _, helmet := range helmetCombinations {
									items[0] = weapon[0]
									items[1] = bodyArmor[0]
									items[2] = legArmor[0]
									items[3] = shield[0]
									items[4] = amulet[0]
									items[5] = boots[0]
									items[6] = ringSet[0]
									items[7] = ringSet[1]
									items[8] = helmet[0]

									f(items)
								}
							}
						}
					}
				}
			}
		}
	}
}

func (c *Character) GetBestGearFor(monsterCode string) ([]oas.SingleItemSchemaItem, error) {
	monster, err := c.GetMonster(monsterCode, true)
	if err != nil {
		return nil, fmt.Errorf("get monster: %w", err)
	}

	items, err := c.bank.Items()
	if err != nil {
		return nil, fmt.Errorf("get bank items: %w", err)
	}

	var weapons, helmets, shields, bodyArmors, legsArmors, boots, rings, amulets []oas.SingleItemSchemaItem
	for _, i := range items {
		if i.Code == "" {
			continue
		}

		item, err := c.GetItem(i.Code, true)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", i.Code, err)
		}

		if item.Level > c.data.Level {
			continue
		}

		switch item.Type {
		case "weapon":
			if item.Code == c.data.WeaponSlot {
				continue
			}
			if item.Subtype == "tool" {
				continue
			}
			weapons = append(weapons, item)
		case "helmet":
			if item.Code == c.data.HelmetSlot {
				continue
			}
			helmets = append(helmets, item)
		case "shield":
			if item.Code == c.data.ShieldSlot {
				continue
			}
			shields = append(shields, item)
		case "body_armor":
			if item.Code == c.data.BodyArmorSlot {
				continue
			}
			bodyArmors = append(bodyArmors, item)
		case "leg_armor":
			if item.Code == c.data.LegArmorSlot {
				continue
			}
			legsArmors = append(legsArmors, item)
		case "boots":
			if item.Code == c.data.BootsSlot {
				continue
			}
			boots = append(boots, item)
		case "ring":
			if item.Code == c.data.Ring1Slot || item.Code == c.data.Ring2Slot {
				continue
			}
			rings = append(rings, item)
		case "amulet":
			if item.Code == c.data.AmuletSlot {
				continue
			}
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
	bestGear := []oas.SingleItemSchemaItem{}

	gearCombinations(func(items []oas.SingleItemSchemaItem) {
		totalEffects := make(map[string]int)

		for _, item := range items {
			for _, e := range item.Effects {
				totalEffects[e.Name] += e.Value
			}
		}

		score := float64(totalEffects["attack_earth"])*(1-float64(monster.ResEarth)/100)*(1+float64(totalEffects["dmg_earth"])/100) +
			float64(totalEffects["attack_water"])*(1-float64(monster.ResWater)/100)*(1+float64(totalEffects["dmg_water"])/100) +
			float64(totalEffects["attack_fire"])*(1-float64(monster.ResFire)/100)*(1+float64(totalEffects["dmg_fire"])/100) +
			float64(totalEffects["attack_air"])*(1-float64(monster.ResAir)/100)*(1+float64(totalEffects["dmg_air"])/100) +
			float64(totalEffects["haste"])
		if score > bestScore {
			bestScore = score
			bestGear = slices.Clone(items)
		}
	}, weaponCombinations, bodyArmorCombinations, legsArmorCombinations,
		shieldCombinations, amuletCombinations, bootsCombinations, ringsCombinations, helmetCombinations)

	bestGearCodes := make([]string, 0, len(bestGear))
	for _, gear := range bestGear {
		bestGearCodes = append(bestGearCodes, gear.Code)
	}
	c.Logger().With(slog.Any("gear", bestGearCodes), slog.Float64("effective_damage", bestScore)).Info("selected best gear for monster: " + monsterCode)

	return bestGear, nil
}

func (c *Character) GetBestGatherGearFor(resourceCode string) ([]oas.SingleItemSchemaItem, error) {
	resource, err := c.GetResource(resourceCode, true)
	if err != nil {
		return nil, err
	}

	items, err := c.bank.Items()
	if err != nil {
		return nil, fmt.Errorf("get bank items: %w", err)
	}

	var tools []oas.SingleItemSchemaItem
	for _, item := range items {
		if item.Subtype == "tool" && item.Level <= c.data.Level {
			tools = append(tools, item.SingleItemSchemaItem)
		}
	}

	if c.data.WeaponSlot != "" {
		item, err := c.GetItem(c.data.WeaponSlot, true)
		if err != nil {
			return nil, fmt.Errorf("get item (current weapon): %w", err)
		}

		if item.Subtype == "tool" {
			tools = append(tools, item)
		}
	}

	slices.SortFunc(tools, func(a, b oas.SingleItemSchemaItem) int {
		aEffect, bEffect := 0, 0

		for _, effect := range a.Effects {
			if effect.Name == string(resource.Skill) {
				aEffect = effect.Value
				break
			}
		}

		for _, effect := range b.Effects {
			if effect.Name == string(resource.Skill) {
				bEffect = effect.Value
				break
			}
		}

		return aEffect - bEffect
	})
	c.Logger().With(slog.Any("gear", []string{tools[0].Code})).Info("selected best gear for resource: " + resource.Code)

	return []oas.SingleItemSchemaItem{tools[0]}, nil
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
	res, err := c.cli.GetBankItemsMyBankItemsGet(context.Background(), oas.GetBankItemsMyBankItemsGetParams{ItemCode: oas.NewOptString(code)})
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
