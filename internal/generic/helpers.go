package generic

import (
	"fmt"
	"strings"

	api "github.com/Sinketsu/artifactsmmo/gen/oas"
	combinations "github.com/mxschmitt/golang-combinations"
)

type ICharacter interface {
	GetMonster(string) (api.MonsterSchema, error)
	Data() api.CharacterSchema
	GetItem(code string) (api.SingleItemSchemaItem, error)
	Log(msg ...any)
}

func getBestGearFor(c ICharacter, monsterCode string) ([][]api.SingleItemSchemaItem, error) {
	monster, err := c.GetMonster(monsterCode)
	if err != nil {
		return nil, fmt.Errorf("get monster: %w", err)
	}

	var weapons, helmets, shields, bodyArmors, legsArmors, boots, rings, amulets []api.SingleItemSchemaItem
	for _, i := range c.Data().Inventory {
		if i.Code == "" {
			continue
		}

		item, err := c.GetItem(i.Code)
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
		item, err := c.GetItem(c.Data().WeaponSlot)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().WeaponSlot, err)
		}
		weapons = append(weapons, item)
	}
	if c.Data().HelmetSlot != "" {
		item, err := c.GetItem(c.Data().HelmetSlot)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().HelmetSlot, err)
		}
		helmets = append(helmets, item)
	}
	if c.Data().BodyArmorSlot != "" {
		item, err := c.GetItem(c.Data().BodyArmorSlot)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().BodyArmorSlot, err)
		}
		bodyArmors = append(bodyArmors, item)
	}
	if c.Data().LegArmorSlot != "" {
		item, err := c.GetItem(c.Data().LegArmorSlot)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().LegArmorSlot, err)
		}
		legsArmors = append(legsArmors, item)
	}
	if c.Data().BootsSlot != "" {
		item, err := c.GetItem(c.Data().BootsSlot)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().BootsSlot, err)
		}
		boots = append(boots, item)
	}
	if c.Data().ShieldSlot != "" {
		item, err := c.GetItem(c.Data().ShieldSlot)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().ShieldSlot, err)
		}
		shields = append(shields, item)
	}
	if c.Data().Ring1Slot != "" {
		item, err := c.GetItem(c.Data().Ring1Slot)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().Ring1Slot, err)
		}
		rings = append(rings, item)
	}
	if c.Data().Ring2Slot != "" {
		item, err := c.GetItem(c.Data().Ring2Slot)
		if err != nil {
			return nil, fmt.Errorf("get item %s: %w", c.Data().Ring2Slot, err)
		}
		rings = append(rings, item)
	}
	if c.Data().AmuletSlot != "" {
		item, err := c.GetItem(c.Data().AmuletSlot)
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
	bestGear := [][]api.SingleItemSchemaItem{}
	for _, weapon := range weaponCombinations {
		for _, bodyArmor := range bodyArmorCombinations {
			for _, legsArmor := range legsArmorCombinations {
				for _, shield := range shieldCombinations {
					for _, amulet := range amuletCombinations {
						for _, boots := range bootsCombinations {
							for _, rings := range ringsCombinations {
								for _, helmet := range helmetCombinations {
									totalEffects := make(map[string]int)

									for _, items := range [][]api.SingleItemSchemaItem{weapon, bodyArmor, legsArmor, shield, amulet, boots, rings, helmet} {
										for _, item := range items {
											for _, e := range item.Effects {
												totalEffects[e.Name] += e.Value
											}
										}
									}

									score := float64(totalEffects["attack_earth"])*(1-float64(monster.ResEarth)/100)*(1+float64(totalEffects["dmg_earth"])/100) +
										float64(totalEffects["attack_water"])*(1-float64(monster.ResWater)/100)*(1+float64(totalEffects["dmg_water"])/100) +
										float64(totalEffects["attack_fire"])*(1-float64(monster.ResFire)/100)*(1+float64(totalEffects["dmg_fire"])/100) +
										float64(totalEffects["attack_air"])*(1-float64(monster.ResAir)/100)*(1+float64(totalEffects["dmg_air"])/100)
									if score > bestScore {
										bestScore = score
										bestGear = [][]api.SingleItemSchemaItem{weapon, bodyArmor, legsArmor, shield, amulet, boots, rings, helmet}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	bestGearCodes := make([]string, len(bestGear))
	for _, gears := range bestGear {
		for _, gear := range gears {
			bestGearCodes = append(bestGearCodes, gear.Code)
		}
	}
	c.Log("found best gear with effective dmg:", bestScore, "[", strings.Join(bestGearCodes, ", "), "]")

	return bestGear, nil
}
