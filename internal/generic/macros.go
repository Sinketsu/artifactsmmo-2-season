package generic

import (
	"fmt"
	"math"
)

func (c *Character) MacroSellAll(codes ...string) error {
	err := c.Move(5, 1) // Grand Exchange
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	for _, code := range codes {
		inInventory := c.InInventory(code)
		if inInventory == 0 {
			continue
		}

		ge, err := c.GetGEItem(code)
		if err != nil {
			return fmt.Errorf("get ge item: %w", err)
		}

		for inInventory > 0 {
			sellCount := int(math.Min(float64(inInventory), float64(ge.MaxQuantity)))

			gold, err := c.Sell(code, sellCount, ge.SellPrice.Value)
			if err != nil {
				return fmt.Errorf("sell: %w", err)
			}

			c.Log("sold", sellCount, code, "Earned", gold, "gold")
			inInventory -= sellCount
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

func (c *Character) MacroGather(code string, sell ...string) error {
	if c.gatherData.Code != code {
		tiles, err := c.FindOnMap(code)
		if err != nil {
			return fmt.Errorf("find on map: %w", err)
		}

		if len(tiles) == 0 {
			return fmt.Errorf("find on map: not found")
		}

		c.gatherData.X = tiles[0].X
		c.gatherData.Y = tiles[0].Y
		c.gatherData.Code = code

		c.Log("found", code, "spot on (", c.gatherData.X, ",", c.gatherData.Y, ")")
	}

	err := c.Move(c.gatherData.X, c.gatherData.Y)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	result, err := c.Gather()
	if err != nil {
		return fmt.Errorf("gather: %w", err)
	}

	c.Log("got", result.Xp, "XP")
	for _, item := range result.Items {
		c.Log("got", item.Quantity, item.Code)
	}

	return nil
}

func (c *Character) MacroFight(monster string, sell ...string) error {
	if c.fightData.Monster != monster {
		tiles, err := c.FindOnMap(monster)
		if err != nil {
			return fmt.Errorf("find on map: %w", err)
		}

		if len(tiles) == 0 {
			return fmt.Errorf("find on map: not found")
		}

		c.fightData.X = tiles[0].X
		c.fightData.Y = tiles[0].Y
		c.fightData.Monster = monster

		c.Log("found", monster, "spot on (", c.fightData.X, ",", c.fightData.Y, ")")
	}

	err := c.Move(c.fightData.X, c.fightData.Y)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	result, err := c.Fight()
	if err != nil {
		return fmt.Errorf("fight: %w", err)
	}

	c.Log("got", result.Xp, "XP")
	for _, item := range result.Drops {
		c.Log("got", item.Quantity, item.Code)
	}

	return nil
}

func (c *Character) MacroCheckCraftResources(code string) (int, error) {
	item, err := c.GetItem(code)
	if err != nil {
		return 0, fmt.Errorf("geet item: %w", err)
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

		c.Log("found", code, "craft spot on (", c.fightData.X, ",", c.fightData.Y, ")")
	}

	err := c.Move(c.craftData.X, c.craftData.Y)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	result, err := c.Craft(code, quantity)
	if err != nil {
		return fmt.Errorf("craft: %w", err)
	}

	c.Log("got", result.Xp, "XP by craft", result.Items[0].Quantity, result.Items[0].Code)
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

		if err := c.Move(tiles[0].X, tiles[1].Y); err != nil {
			return fmt.Errorf("move: %w", err)
		}

		result, err := c.Recycle(code, inInventory)
		if err != nil {
			return fmt.Errorf("recycle: %w", err)
		}

		for _, item := range result.Items {
			c.Log("got", item.Quantity, item.Code, "via recycling")
		}
	}

	return nil
}
