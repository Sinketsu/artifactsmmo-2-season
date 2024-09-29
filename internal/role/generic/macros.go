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

	totalSold := 0

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
			totalSold += sellCount
		}
	}

	if totalSold == 0 {
		return fmt.Errorf("inventory not contains all of codes: %v", codes)
	}

	return nil
}

func (c *Character) MacroGather(code string, sell ...string) error {
	if c.InventoryItemCount() == c.Data().InventoryMaxItems {
		err := c.MacroSellAll(sell...)
		if err != nil {
			return fmt.Errorf("sell: %w", err)
		}
	}

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
	if c.InventoryItemCount() == c.Data().InventoryMaxItems {
		err := c.MacroSellAll(sell...)
		if err != nil {
			return fmt.Errorf("sell: %w", err)
		}
	}

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
