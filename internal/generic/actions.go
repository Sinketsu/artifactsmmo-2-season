package generic

import (
	"context"
	"fmt"
	"time"
	"unsafe"

	oas "github.com/Sinketsu/artifactsmmo/gen/oas"
)

func (c *Character) Gather() (oas.SkillDataSchemaDetails, error) {
	requestCount.Inc()

	res, err := c.cli.ActionGatheringMyNameActionGatheringPost(context.Background(), oas.ActionGatheringMyNameActionGatheringPostParams{Name: c.name})
	if err != nil {
		return oas.SkillDataSchemaDetails{}, err
	}

	switch v := res.(type) {
	case *oas.SkillResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return v.Data.Details, c.updateData(unsafe.Pointer(&v.Data.Character))
	case *oas.ActionGatheringMyNameActionGatheringPostCode486:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("already gathering...")
	case *oas.ActionGatheringMyNameActionGatheringPostCode493:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("character skill level is too low")
	case *oas.ActionGatheringMyNameActionGatheringPostCode497:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("inventory is full")
	case *oas.ActionGatheringMyNameActionGatheringPostCode498:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("character not found")
	case *oas.ActionGatheringMyNameActionGatheringPostCode499:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("cooldown")
	case *oas.ActionGatheringMyNameActionGatheringPostCode598:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("not on required map tile")
	default:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Move(x, y int) error {
	if c.data.X == x && c.data.Y == y {
		return nil
	}

	requestCount.Inc()

	res, err := c.cli.ActionMoveMyNameActionMovePost(context.Background(), &oas.DestinationSchema{X: x, Y: y}, oas.ActionMoveMyNameActionMovePostParams{Name: c.name})
	if err != nil {
		return err
	}

	switch v := res.(type) {
	case *oas.CharacterMovementResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return c.updateData(unsafe.Pointer(&v.Data.Character))
	case *oas.ActionMoveMyNameActionMovePostNotFound:
		return fmt.Errorf("map not found")
	case *oas.ActionMoveMyNameActionMovePostCode486:
		return fmt.Errorf("already moving...")
	case *oas.ActionMoveMyNameActionMovePostCode490:
		// character already at point
		return nil
	case *oas.ActionMoveMyNameActionMovePostCode498:
		return fmt.Errorf("character not found")
	case *oas.ActionMoveMyNameActionMovePostCode499:
		return fmt.Errorf("cooldown")
	default:
		return fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Sell(code string, quantity int, price int) (int, error) {
	requestCount.Inc()

	res, err := c.cli.ActionGeSellItemMyNameActionGeSellPost(context.Background(), &oas.GETransactionItemSchema{Code: code, Quantity: quantity, Price: price}, oas.ActionGeSellItemMyNameActionGeSellPostParams{Name: c.name})
	if err != nil {
		return 0, err
	}

	switch v := res.(type) {
	case *oas.GETransactionResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return v.Data.Transaction.TotalPrice, c.updateData(unsafe.Pointer(&v.Data.Character))
	case *oas.ActionGeSellItemMyNameActionGeSellPostNotFound:
		return 0, fmt.Errorf("item not found")
	case *oas.ActionGeSellItemMyNameActionGeSellPostCode478:
		return 0, fmt.Errorf("missing item or insufficient quantity")
	case *oas.ActionGeSellItemMyNameActionGeSellPostCode479:
		return 0, fmt.Errorf("too many items to sell - bigger than limit")
	case *oas.ActionGeSellItemMyNameActionGeSellPostCode482:
		return 0, fmt.Errorf("no item at this price")
	case *oas.ActionGeSellItemMyNameActionGeSellPostCode483:
		return 0, fmt.Errorf("transaction is already in progress on this item by a another character")
	case *oas.ActionGeSellItemMyNameActionGeSellPostCode486:
		return 0, fmt.Errorf("action is already in progress by your character")
	case *oas.ActionGeSellItemMyNameActionGeSellPostCode498:
		return 0, fmt.Errorf("character not found")
	case *oas.ActionGeSellItemMyNameActionGeSellPostCode499:
		return 0, fmt.Errorf("cooldown")
	case *oas.ActionGeSellItemMyNameActionGeSellPostCode598:
		return 0, fmt.Errorf("GE not at this map tile")
	default:
		return 0, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) GetGEItem(code string) (oas.GEItemSchema, error) {
	requestCount.Inc()

	res, err := c.cli.GetGeItemGeCodeGet(context.Background(), oas.GetGeItemGeCodeGetParams{Code: code})
	if err != nil {
		return oas.GEItemSchema{}, err
	}

	switch v := res.(type) {
	case *oas.GEItemResponseSchema:
		return v.Data, nil
	case *oas.GetGeItemGeCodeGetNotFound:
		return oas.GEItemSchema{}, fmt.Errorf("item not found")
	default:
		return oas.GEItemSchema{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) GetItem(code string, cachable bool) (oas.SingleItemSchemaItem, error) {
	if cachable {
		if item, ok := c.itemsCache[code]; ok {
			return item, nil
		}
	}

	requestCount.Inc()

	res, err := c.cli.GetItemItemsCodeGet(context.Background(), oas.GetItemItemsCodeGetParams{Code: code})
	if err != nil {
		return oas.SingleItemSchemaItem{}, err
	}

	switch v := res.(type) {
	case *oas.ItemResponseSchema:
		if cachable {
			c.itemsCache[code] = v.Data.Item
		}

		return v.Data.Item, nil
	case *oas.GetItemItemsCodeGetNotFound:
		return oas.SingleItemSchemaItem{}, fmt.Errorf("item not found")
	default:
		return oas.SingleItemSchemaItem{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Fight() (oas.CharacterFightDataSchemaFight, error) {
	requestCount.Inc()

	res, err := c.cli.ActionFightMyNameActionFightPost(context.Background(), oas.ActionFightMyNameActionFightPostParams{Name: c.name})
	if err != nil {
		return oas.CharacterFightDataSchemaFight{}, err
	}

	switch v := res.(type) {
	case *oas.CharacterFightResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		if v.Data.Fight.Result == oas.CharacterFightDataSchemaFightResultLose {
			return oas.CharacterFightDataSchemaFight{}, fmt.Errorf("loose battle")
		}

		return v.Data.Fight, c.updateData(unsafe.Pointer(&v.Data.Character))
	case *oas.ActionFightMyNameActionFightPostCode486:
		return oas.CharacterFightDataSchemaFight{}, fmt.Errorf("action is already in progress by your character")
	case *oas.ActionFightMyNameActionFightPostCode497:
		return oas.CharacterFightDataSchemaFight{}, fmt.Errorf("inventory is full")
	case *oas.ActionFightMyNameActionFightPostCode498:
		return oas.CharacterFightDataSchemaFight{}, fmt.Errorf("character not found")
	case *oas.ActionFightMyNameActionFightPostCode499:
		return oas.CharacterFightDataSchemaFight{}, fmt.Errorf("cooldown")
	case *oas.ActionFightMyNameActionFightPostCode598:
		return oas.CharacterFightDataSchemaFight{}, fmt.Errorf("monster is not at this map tile")
	default:
		return oas.CharacterFightDataSchemaFight{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Craft(code string, quantity int) (oas.SkillDataSchemaDetails, error) {
	requestCount.Inc()

	res, err := c.cli.ActionCraftingMyNameActionCraftingPost(context.Background(), &oas.CraftingSchema{Code: code, Quantity: oas.NewOptInt(quantity)}, oas.ActionCraftingMyNameActionCraftingPostParams{Name: c.name})
	if err != nil {
		return oas.SkillDataSchemaDetails{}, err
	}

	switch v := res.(type) {
	case *oas.SkillResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return v.Data.Details, c.updateData(unsafe.Pointer(&v.Data.Character))
	case *oas.ActionCraftingMyNameActionCraftingPostNotFound:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("craft not found")
	case *oas.ActionCraftingMyNameActionCraftingPostCode478:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("missing item or insufficient quantity")
	case *oas.ActionCraftingMyNameActionCraftingPostCode486:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("action is already in progress by your character")
	case *oas.ActionCraftingMyNameActionCraftingPostCode493:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("skill level is too low")
	case *oas.ActionCraftingMyNameActionCraftingPostCode497:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("inventory is full")
	case *oas.ActionCraftingMyNameActionCraftingPostCode498:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("character not found")
	case *oas.ActionCraftingMyNameActionCraftingPostCode499:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("cooldown")
	case *oas.ActionCraftingMyNameActionCraftingPostCode598:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("workshop not at this map tile")
	default:
		return oas.SkillDataSchemaDetails{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Withdraw(code string, quantity int) error {
	requestCount.Inc()

	res, err := c.cli.ActionWithdrawBankMyNameActionBankWithdrawPost(context.Background(), &oas.SimpleItemSchema{Code: code, Quantity: quantity}, oas.ActionWithdrawBankMyNameActionBankWithdrawPostParams{Name: c.name})
	if err != nil {
		return err
	}

	switch v := res.(type) {
	case *oas.BankItemTransactionResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return c.updateData(unsafe.Pointer(&v.Data.Character))
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostNotFound:
		return fmt.Errorf("item not found")
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostCode461:
		return fmt.Errorf("transaction is already in progress with this item/your golds in your bank")
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostCode478:
		return fmt.Errorf("missing item or insufficient quantity")
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostCode486:
		return fmt.Errorf("action is already in progress by your character")
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostCode497:
		return fmt.Errorf("inventory is full")
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostCode498:
		return fmt.Errorf("character not found")
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostCode499:
		return fmt.Errorf("cooldown")
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostCode598:
		return fmt.Errorf("bank not at this map tile")
	default:
		return fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Deposit(code string, quantity int) error {
	requestCount.Inc()

	res, err := c.cli.ActionDepositBankMyNameActionBankDepositPost(context.Background(), &oas.SimpleItemSchema{Code: code, Quantity: quantity}, oas.ActionDepositBankMyNameActionBankDepositPostParams{Name: c.name})
	if err != nil {
		return err
	}

	switch v := res.(type) {
	case *oas.BankItemTransactionResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return c.updateData(unsafe.Pointer(&v.Data.Character))
	case *oas.ActionDepositBankMyNameActionBankDepositPostNotFound:
		return fmt.Errorf("item not found")
	case *oas.ActionDepositBankMyNameActionBankDepositPostCode461:
		return fmt.Errorf("transaction is already in progress with this item/your golds in your bank")
	case *oas.ActionDepositBankMyNameActionBankDepositPostCode462:
		return fmt.Errorf("bank is full")
	case *oas.ActionDepositBankMyNameActionBankDepositPostCode478:
		return fmt.Errorf("missing item or insufficient quantity")
	case *oas.ActionDepositBankMyNameActionBankDepositPostCode486:
		return fmt.Errorf("action is already in progress by your character")
	case *oas.ActionDepositBankMyNameActionBankDepositPostCode498:
		return fmt.Errorf("character not found")
	case *oas.ActionDepositBankMyNameActionBankDepositPostCode499:
		return fmt.Errorf("cooldown")
	case *oas.ActionDepositBankMyNameActionBankDepositPostCode598:
		return fmt.Errorf("bank not at this map tile")
	default:
		return fmt.Errorf("unknown answer type")
	}
}

func (c *Character) DepositGold(quantity int) error {
	requestCount.Inc()

	res, err := c.cli.ActionDepositBankGoldMyNameActionBankDepositGoldPost(context.Background(), &oas.DepositWithdrawGoldSchema{Quantity: quantity}, oas.ActionDepositBankGoldMyNameActionBankDepositGoldPostParams{Name: c.name})
	if err != nil {
		return err
	}

	switch v := res.(type) {
	case *oas.BankGoldTransactionResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return c.updateData(unsafe.Pointer(&v.Data.Character))
	case *oas.ActionDepositBankGoldMyNameActionBankDepositGoldPostCode461:
		return fmt.Errorf("transaction is already in progress with this item/your golds in your bank")
	case *oas.ActionDepositBankGoldMyNameActionBankDepositGoldPostCode486:
		return fmt.Errorf("action is already in progress by your character")
	case *oas.ActionDepositBankGoldMyNameActionBankDepositGoldPostCode492:
		return fmt.Errorf("insufficient gold")
	case *oas.ActionDepositBankGoldMyNameActionBankDepositGoldPostCode498:
		return fmt.Errorf("character not found")
	case *oas.ActionDepositBankGoldMyNameActionBankDepositGoldPostCode499:
		return fmt.Errorf("cooldown")
	case *oas.ActionDepositBankGoldMyNameActionBankDepositGoldPostCode598:
		return fmt.Errorf("bank not at this map tile")
	default:
		return fmt.Errorf("unknown answer type")
	}
}

func (c *Character) FindOnMap(code string, cachable bool) (oas.MapSchema, error) {
	if cachable {
		if tile, ok := c.mapsCache[code]; ok {
			return tile, nil
		}
	}

	requestCount.Inc()

	res, err := c.cli.GetAllMapsMapsGet(context.Background(), oas.GetAllMapsMapsGetParams{ContentCode: oas.NewOptString(code)})
	if err != nil {
		return oas.MapSchema{}, err
	}

	if len(res.Data) == 0 {
		return oas.MapSchema{}, fmt.Errorf("not found")
	}

	if cachable {
		c.mapsCache[code] = res.Data[0]
	}

	return res.Data[0], nil
}

func (c *Character) GetResource(code string, cachable bool) (oas.ResourceSchema, error) {
	if cachable {
		if resource, ok := c.resourceCache[code]; ok {
			return resource, nil
		}
	}

	requestCount.Inc()

	res, err := c.cli.GetResourceResourcesCodeGet(context.Background(), oas.GetResourceResourcesCodeGetParams{Code: code})
	if err != nil {
		return oas.ResourceSchema{}, err
	}

	switch v := res.(type) {
	case *oas.ResourceResponseSchema:
		if cachable {
			c.resourceCache[code] = v.Data
		}

		return v.Data, nil
	case *oas.GetResourceResourcesCodeGetNotFound:
		return oas.ResourceSchema{}, fmt.Errorf("resource not found")
	default:
		return oas.ResourceSchema{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) GetMonster(code string, cachable bool) (oas.MonsterSchema, error) {
	if cachable {
		if monster, ok := c.monsterCache[code]; ok {
			return monster, nil
		}
	}

	requestCount.Inc()

	res, err := c.cli.GetMonsterMonstersCodeGet(context.Background(), oas.GetMonsterMonstersCodeGetParams{Code: code})
	if err != nil {
		return oas.MonsterSchema{}, err
	}

	switch v := res.(type) {
	case *oas.MonsterResponseSchema:
		if cachable {
			c.monsterCache[code] = v.Data
		}

		return v.Data, nil
	case *oas.GetMonsterMonstersCodeGetNotFound:
		return oas.MonsterSchema{}, fmt.Errorf("monster not found")
	default:
		return oas.MonsterSchema{}, fmt.Errorf("unknown answer type: %v", v)
	}
}

func (c *Character) Recycle(code string, quantity int) (oas.RecyclingDataSchemaDetails, error) {
	requestCount.Inc()

	res, err := c.cli.ActionRecyclingMyNameActionRecyclingPost(context.Background(), &oas.RecyclingSchema{Code: code, Quantity: oas.NewOptInt(quantity)}, oas.ActionRecyclingMyNameActionRecyclingPostParams{Name: c.name})
	if err != nil {
		return oas.RecyclingDataSchemaDetails{}, err
	}

	switch v := res.(type) {
	case *oas.RecyclingResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return v.Data.Details, c.updateData(unsafe.Pointer(&v.Data.Character))
	case *oas.ActionRecyclingMyNameActionRecyclingPostNotFound:
		return oas.RecyclingDataSchemaDetails{}, fmt.Errorf("item not found")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode473:
		return oas.RecyclingDataSchemaDetails{}, fmt.Errorf("item cannot be recycled")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode478:
		return oas.RecyclingDataSchemaDetails{}, fmt.Errorf("missing item or insufficient quantity")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode486:
		return oas.RecyclingDataSchemaDetails{}, fmt.Errorf("action is already in progress by your character")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode493:
		return oas.RecyclingDataSchemaDetails{}, fmt.Errorf("skill level is too low")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode497:
		return oas.RecyclingDataSchemaDetails{}, fmt.Errorf("inventory is full")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode498:
		return oas.RecyclingDataSchemaDetails{}, fmt.Errorf("character not found")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode499:
		return oas.RecyclingDataSchemaDetails{}, fmt.Errorf("cooldown")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode598:
		return oas.RecyclingDataSchemaDetails{}, fmt.Errorf("workshop not found on this map")
	default:
		return oas.RecyclingDataSchemaDetails{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Equip(code string, slot string, quantity int) error {
	requestCount.Inc()

	res, err := c.cli.ActionEquipItemMyNameActionEquipPost(context.Background(), &oas.EquipSchema{Code: code, Slot: oas.EquipSchemaSlot(slot), Quantity: oas.NewOptInt(quantity)}, oas.ActionEquipItemMyNameActionEquipPostParams{Name: c.name})
	if err != nil {
		return err
	}

	switch v := res.(type) {
	case *oas.EquipmentResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return c.updateData(unsafe.Pointer(&v.Data.Character))
	case *oas.ActionEquipItemMyNameActionEquipPostNotFound:
		return fmt.Errorf("item not found")
	case *oas.ActionEquipItemMyNameActionEquipPostCode478:
		return fmt.Errorf("missing item or insufficient quantity")
	case *oas.ActionEquipItemMyNameActionEquipPostCode484:
		return fmt.Errorf("character can't equip more than 100 consumables in the same slot")
	case *oas.ActionEquipItemMyNameActionEquipPostCode485:
		return fmt.Errorf("item is already equipped")
	case *oas.ActionEquipItemMyNameActionEquipPostCode486:
		return fmt.Errorf("action is already in progress by your character")
	case *oas.ActionEquipItemMyNameActionEquipPostCode491:
		return fmt.Errorf("slot is not empty")
	case *oas.ActionEquipItemMyNameActionEquipPostCode496:
		return fmt.Errorf("character level is too low")
	case *oas.ActionEquipItemMyNameActionEquipPostCode497:
		return fmt.Errorf("inventory is full")
	case *oas.ActionEquipItemMyNameActionEquipPostCode498:
		return fmt.Errorf("character not found")
	case *oas.ActionEquipItemMyNameActionEquipPostCode499:
		return fmt.Errorf("cooldown")
	default:
		return fmt.Errorf("unknown answer type: %v", v)
	}
}

func (c *Character) UnEquip(code string, slot string, quantity int) error {
	requestCount.Inc()

	res, err := c.cli.ActionUnequipItemMyNameActionUnequipPost(context.Background(), &oas.UnequipSchema{Slot: oas.UnequipSchemaSlot(slot), Quantity: oas.NewOptInt(quantity)}, oas.ActionUnequipItemMyNameActionUnequipPostParams{Name: c.name})
	if err != nil {
		return err
	}

	switch v := res.(type) {
	case *oas.EquipmentResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return c.updateData(unsafe.Pointer(&v.Data.Character))
	case *oas.ActionUnequipItemMyNameActionUnequipPostNotFound:
		return fmt.Errorf("item not found")
	case *oas.ActionUnequipItemMyNameActionUnequipPostCode478:
		return fmt.Errorf("missing item or insufficient quantity")
	case *oas.ActionUnequipItemMyNameActionUnequipPostCode486:
		return fmt.Errorf("action is already in progress by your character")
	case *oas.ActionUnequipItemMyNameActionUnequipPostCode491:
		return fmt.Errorf("slot is empty")
	case *oas.ActionUnequipItemMyNameActionUnequipPostCode497:
		return fmt.Errorf("inventory is full")
	case *oas.ActionUnequipItemMyNameActionUnequipPostCode498:
		return fmt.Errorf("character not found")
	case *oas.ActionUnequipItemMyNameActionUnequipPostCode499:
		return fmt.Errorf("cooldown")
	default:
		return fmt.Errorf("unknown answer type: %v", v)
	}
}

func (c *Character) CompleteTask() (oas.TasksRewardDataSchemaReward, error) {
	requestCount.Inc()

	res, err := c.cli.ActionCompleteTaskMyNameActionTaskCompletePost(context.Background(), oas.ActionCompleteTaskMyNameActionTaskCompletePostParams{Name: c.name})
	if err != nil {
		return oas.TasksRewardDataSchemaReward{}, err
	}

	switch v := res.(type) {
	case *oas.TasksRewardResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return v.Data.Reward, c.updateData(unsafe.Pointer(&v.Data.Character))
	case *oas.ActionCompleteTaskMyNameActionTaskCompletePostCode486:
		return oas.TasksRewardDataSchemaReward{}, fmt.Errorf("action is already in progress by your character")
	case *oas.ActionCompleteTaskMyNameActionTaskCompletePostCode487:
		return oas.TasksRewardDataSchemaReward{}, fmt.Errorf("character has no task")
	case *oas.ActionCompleteTaskMyNameActionTaskCompletePostCode488:
		return oas.TasksRewardDataSchemaReward{}, fmt.Errorf("character has not completed the task")
	case *oas.ActionCompleteTaskMyNameActionTaskCompletePostCode497:
		return oas.TasksRewardDataSchemaReward{}, fmt.Errorf("inventory is full")
	case *oas.ActionCompleteTaskMyNameActionTaskCompletePostCode498:
		return oas.TasksRewardDataSchemaReward{}, fmt.Errorf("character not found")
	case *oas.ActionCompleteTaskMyNameActionTaskCompletePostCode499:
		return oas.TasksRewardDataSchemaReward{}, fmt.Errorf("cooldown")
	case *oas.ActionCompleteTaskMyNameActionTaskCompletePostCode598:
		return oas.TasksRewardDataSchemaReward{}, fmt.Errorf("tasks master is not at this tile map")
	default:
		return oas.TasksRewardDataSchemaReward{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) AcceptNewTask() (oas.TaskDataSchemaTask, error) {
	requestCount.Inc()

	res, err := c.cli.ActionAcceptNewTaskMyNameActionTaskNewPost(context.Background(), oas.ActionAcceptNewTaskMyNameActionTaskNewPostParams{Name: c.name})
	if err != nil {
		return oas.TaskDataSchemaTask{}, err
	}

	switch v := res.(type) {
	case *oas.TaskResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return v.Data.Task, c.updateData(unsafe.Pointer(&v.Data.Character))
	case *oas.ActionAcceptNewTaskMyNameActionTaskNewPostCode486:
		return oas.TaskDataSchemaTask{}, fmt.Errorf("action is already in progress by your character")
	case *oas.ActionAcceptNewTaskMyNameActionTaskNewPostCode489:
		return oas.TaskDataSchemaTask{}, fmt.Errorf("character already has a task")
	case *oas.ActionAcceptNewTaskMyNameActionTaskNewPostCode498:
		return oas.TaskDataSchemaTask{}, fmt.Errorf("character not found")
	case *oas.ActionAcceptNewTaskMyNameActionTaskNewPostCode499:
		return oas.TaskDataSchemaTask{}, fmt.Errorf("cooldown")
	case *oas.ActionAcceptNewTaskMyNameActionTaskNewPostCode598:
		return oas.TaskDataSchemaTask{}, fmt.Errorf("tasks master is not at this tile map")
	default:
		return oas.TaskDataSchemaTask{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) CancelTask() error {
	requestCount.Inc()

	res, err := c.cli.ActionTaskCancelMyNameActionTaskCancelPost(context.Background(), oas.ActionTaskCancelMyNameActionTaskCancelPostParams{Name: c.name})
	if err != nil {
		return err
	}

	switch v := res.(type) {
	case *oas.TaskCancelledResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return c.updateData(unsafe.Pointer(&v.Data.Character))
	case *oas.ActionTaskCancelMyNameActionTaskCancelPostCode478:
		return fmt.Errorf("missing item or insufficient quantity")
	case *oas.ActionTaskCancelMyNameActionTaskCancelPostCode486:
		return fmt.Errorf("action is already in progress by your character")
	case *oas.ActionTaskCancelMyNameActionTaskCancelPostCode498:
		return fmt.Errorf("character not found")
	case *oas.ActionTaskCancelMyNameActionTaskCancelPostCode499:
		return fmt.Errorf("cooldown")
	case *oas.ActionTaskCancelMyNameActionTaskCancelPostCode598:
		return fmt.Errorf("tasks master is not at this tile map")
	default:
		return fmt.Errorf("unknown answer type")
	}
}
