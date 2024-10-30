package generic

import (
	"fmt"
	"log/slog"
	"math"

	oas "github.com/Sinketsu/artifactsmmo/gen/oas"
)

func (c *Character) MacroWithdraw(code string, quantity int) error {
	err := c.Move(4, 1) // Bank
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	if err := c.Withdraw(code, quantity); err != nil {
		return err
	}

	c.Logger().Info(fmt.Sprintf("withdraw: %d %s", quantity, code))
	return nil
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

			price, err := c.Sell(code, sellCount, ge.SellPrice.Value)
			if err != nil {
				return fmt.Errorf("sell: %w", err)
			}
			c.Logger().With(slog.Int("gold", price)).Info(fmt.Sprintf("sold %d %s", sellCount, code))
		}
	}

	return nil
}

func (c *Character) MacroDeposit(code string, quantity int) error {
	err := c.Move(4, 1) // Bank
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	if err := c.Deposit(code, quantity); err != nil {
		return err
	}

	c.Logger().Info(fmt.Sprintf("deposit: %d %s", quantity, code))
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
		c.Logger().Info(fmt.Sprintf("deposit: %d %s", inInventory, code))
	}

	return nil
}

func (c *Character) MacroDepositGold(count int) error {
	err := c.Move(4, 1) // Bank
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	if err := c.DepositGold(count); err != nil {
		return err
	}

	c.Logger().Info(fmt.Sprintf("deposit gold: %d", count))
	return nil
}

func (c *Character) MacroWithdrawGold(count int) error {
	err := c.Move(4, 1) // Bank
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	if err := c.WithdrawGold(count); err != nil {
		return err
	}

	c.Logger().Info(fmt.Sprintf("withdraw gold: %d", count))
	return nil
}

func (c *Character) MacroGather(resource string, cachable bool) error {
	tile, err := c.FindOnMap(resource, cachable)
	if err != nil {
		return fmt.Errorf("find on map: %w", err)
	}

	err = c.Move(tile.X, tile.Y)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	gathered, err := c.Gather()
	if err != nil {
		return fmt.Errorf("gather: %w", err)
	}

	c.Logger().With(
		slog.Int("xp", gathered.Xp),
		slog.Any("items", gathered.Items),
	).Info("gather: " + resource)
	return nil
}

func (c *Character) MacroFight(monster string, cachable bool) error {
	tile, err := c.FindOnMap(monster, cachable)
	if err != nil {
		return fmt.Errorf("find on map: %w", err)
	}

	err = c.Move(tile.X, tile.Y)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	result, err := c.Fight()
	if err != nil {
		return fmt.Errorf("fight: %w", err)
	}
	c.Logger().With(
		slog.Int("xp", result.Xp),
		slog.Int("gold", result.Gold),
		slog.Any("items", result.Drops),
		slog.Int("turns", result.Turns),
	).Info("fight: " + monster)

	return nil
}

func (c *Character) MacroCheckCraftResources(code string) (int, error) {
	item, err := c.GetItem(code, true)
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
	item, err := c.GetItem(code, true)
	if err != nil {
		return fmt.Errorf("get item: %w", err)
	}

	workshop, err := c.FindOnMap(string(item.Craft.Value.CraftSchema.Skill.Value), true)
	if err != nil {
		return fmt.Errorf("find on map: %w", err)
	}

	err = c.Move(workshop.X, workshop.Y)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	crafted, err := c.Craft(code, quantity)
	if err != nil {
		return fmt.Errorf("craft: %w", err)
	}
	c.Logger().With(slog.Int("xp", crafted.Xp), slog.Any("items", crafted.Items)).Info(fmt.Sprintf("craft: %d %s", quantity, code))

	return nil
}

func (c *Character) MacroBuy(code string, quantity int, price int) error {
	err := c.Move(5, 1) // Grand Exchange
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	gold, err := c.Buy(code, quantity, price)
	if err != nil {
		return fmt.Errorf("buy: %w", err)
	}

	c.Logger().With(slog.Int("total", gold)).Info(fmt.Sprintf("buy: %d %s", quantity, code))
	return nil
}

func (c *Character) MacroRecycleAll(codes ...string) error {
	for _, code := range codes {
		inInventory := c.InInventory(code)
		if inInventory == 0 {
			continue
		}

		item, err := c.GetItem(code, true)
		if err != nil {
			return fmt.Errorf("get item: %w", err)
		}

		workshop, err := c.FindOnMap(string(item.Craft.Value.CraftSchema.Skill.Value), true)
		if err != nil {
			return fmt.Errorf("find on map: %w", err)
		}

		if err := c.Move(workshop.X, workshop.Y); err != nil {
			return fmt.Errorf("move: %w", err)
		}

		items, err := c.Recycle(code, inInventory)
		if err != nil {
			return fmt.Errorf("recycle: %w", err)
		}
		c.Logger().With(slog.Any("items", items)).Info(fmt.Sprintf("recycle: %d %s", inInventory, code))
	}

	return nil
}

func (c *Character) MacroWear(items []oas.SingleItemSchemaItem) error {
	ringCount := 0

	err := c.Move(4, 1) // Bank
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	for _, item := range items {
		var current string
		slot := item.Type

		switch item.Type {
		case "weapon":
			current = c.data.WeaponSlot
		case "helmet":
			current = c.data.HelmetSlot
		case "body_armor":
			current = c.data.BodyArmorSlot
		case "leg_armor":
			current = c.data.LegArmorSlot
		case "shield":
			current = c.data.ShieldSlot
		case "boots":
			current = c.data.BootsSlot
		case "amulet":
			current = c.data.AmuletSlot
		case "ring":
			switch ringCount {
			case 0:
				current = c.data.Ring1Slot
				slot = "ring1"
				ringCount++
			case 1:
				current = c.data.Ring2Slot
				slot = "ring2"
				ringCount++
			}
		default:
			return fmt.Errorf("unknown type: %s", item.Type)
		}

		if current == item.Code {
			continue
		}

		if current != "" {
			if err := c.UnEquip(current, slot, 1); err != nil {
				return fmt.Errorf("unequip: %w", err)
			}
			c.Logger().Info("unequip: " + current)

			if err := c.Deposit(current, 1); err != nil {
				return fmt.Errorf("deposit: %w", err)
			}
			c.Logger().Info(fmt.Sprintf("deposit: %d %s", 1, current))
		}

		if err := c.Withdraw(item.Code, 1); err != nil {
			return fmt.Errorf("withdraw: %w", err)
		}
		c.Logger().Info(fmt.Sprintf("withdraw: %d %s", 1, item.Code))

		if err := c.Equip(item.Code, slot, 1); err != nil {
			return err
		}
		c.Logger().Info("equip: " + item.Code)
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

	c.Logger().With(slog.Int(reward.Code, reward.Quantity)).Info("complete task")
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

	c.Logger().With(slog.Int(task.Code, task.Total)).Info("accept new task")
	return nil
}

func (c *Character) MacroCancelMonsterTask() error {
	err := c.Move(1, 2) // Task master monsters
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	task := c.Data().Task

	err = c.CancelTask()
	if err != nil {
		return fmt.Errorf("cancel task: %w", err)
	}

	c.Logger().With(slog.String("task", task)).Warn("cancel task")
	return nil
}
