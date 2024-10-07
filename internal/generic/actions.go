package generic

import (
	"context"
	"fmt"
	"time"
	"unsafe"

	api "github.com/Sinketsu/artifactsmmo/gen/oas"
)

func (c *Character) Gather() (api.SkillDataSchemaDetails, error) {
	res, err := c.cli.ActionGatheringMyNameActionGatheringPost(context.Background(), api.ActionGatheringMyNameActionGatheringPostParams{Name: c.name})
	if err != nil {
		return api.SkillDataSchemaDetails{}, err
	}

	switch v := res.(type) {
	case *api.SkillResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return v.Data.Details, c.updateData(unsafe.Pointer(&v.Data.Character))
	case *api.ActionGatheringMyNameActionGatheringPostCode486:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("already gathering...")
	case *api.ActionGatheringMyNameActionGatheringPostCode493:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("character skill level is too low")
	case *api.ActionGatheringMyNameActionGatheringPostCode497:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("inventory is full")
	case *api.ActionGatheringMyNameActionGatheringPostCode498:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("character not found")
	case *api.ActionGatheringMyNameActionGatheringPostCode499:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("cooldown")
	case *api.ActionGatheringMyNameActionGatheringPostCode598:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("not on required map tile")
	default:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Move(x, y int) error {
	if c.data.X == x && c.data.Y == y {
		return nil
	}

	res, err := c.cli.ActionMoveMyNameActionMovePost(context.Background(), &api.DestinationSchema{X: x, Y: y}, api.ActionMoveMyNameActionMovePostParams{Name: c.name})
	if err != nil {
		return err
	}

	switch v := res.(type) {
	case *api.CharacterMovementResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return c.updateData(unsafe.Pointer(&v.Data.Character))
	case *api.ActionMoveMyNameActionMovePostNotFound:
		return fmt.Errorf("map not found")
	case *api.ActionMoveMyNameActionMovePostCode486:
		return fmt.Errorf("already moving...")
	case *api.ActionMoveMyNameActionMovePostCode490:
		return fmt.Errorf("character already at point")
	case *api.ActionMoveMyNameActionMovePostCode498:
		return fmt.Errorf("character not found")
	case *api.ActionMoveMyNameActionMovePostCode499:
		return fmt.Errorf("cooldown")
	default:
		return fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Sell(code string, quantity int, price int) (int, error) {
	res, err := c.cli.ActionGeSellItemMyNameActionGeSellPost(context.Background(), &api.GETransactionItemSchema{Code: code, Quantity: quantity, Price: price}, api.ActionGeSellItemMyNameActionGeSellPostParams{Name: c.name})
	if err != nil {
		return 0, err
	}

	switch v := res.(type) {
	case *api.GETransactionResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return v.Data.Transaction.TotalPrice, c.updateData(unsafe.Pointer(&v.Data.Character))
	case *api.ActionGeSellItemMyNameActionGeSellPostNotFound:
		return 0, fmt.Errorf("item not found")
	case *api.ActionGeSellItemMyNameActionGeSellPostCode478:
		return 0, fmt.Errorf("missing item or insufficient quantity")
	case *api.ActionGeSellItemMyNameActionGeSellPostCode479:
		return 0, fmt.Errorf("too many items to sell - bigger than limit")
	case *api.ActionGeSellItemMyNameActionGeSellPostCode482:
		return 0, fmt.Errorf("no item at this price")
	case *api.ActionGeSellItemMyNameActionGeSellPostCode483:
		return 0, fmt.Errorf("transaction is already in progress on this item by a another character")
	case *api.ActionGeSellItemMyNameActionGeSellPostCode486:
		return 0, fmt.Errorf("action is already in progress by your character")
	case *api.ActionGeSellItemMyNameActionGeSellPostCode498:
		return 0, fmt.Errorf("character not found")
	case *api.ActionGeSellItemMyNameActionGeSellPostCode499:
		return 0, fmt.Errorf("cooldown")
	case *api.ActionGeSellItemMyNameActionGeSellPostCode598:
		return 0, fmt.Errorf("GE not at this map tile")
	default:
		return 0, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) GetGEItem(code string) (api.GEItemSchema, error) {
	res, err := c.cli.GetGeItemGeCodeGet(context.Background(), api.GetGeItemGeCodeGetParams{Code: code})
	if err != nil {
		return api.GEItemSchema{}, err
	}

	switch v := res.(type) {
	case *api.GEItemResponseSchema:
		return v.Data, nil
	case *api.GetGeItemGeCodeGetNotFound:
		return api.GEItemSchema{}, fmt.Errorf("item not found")
	default:
		return api.GEItemSchema{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) GetItem(code string) (api.SingleItemSchemaItem, error) {
	res, err := c.cli.GetItemItemsCodeGet(context.Background(), api.GetItemItemsCodeGetParams{Code: code})
	if err != nil {
		return api.SingleItemSchemaItem{}, err
	}

	switch v := res.(type) {
	case *api.ItemResponseSchema:
		return v.Data.Item, nil
	case *api.GetItemItemsCodeGetNotFound:
		return api.SingleItemSchemaItem{}, fmt.Errorf("item not found")
	default:
		return api.SingleItemSchemaItem{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Fight() (api.CharacterFightDataSchemaFight, error) {
	res, err := c.cli.ActionFightMyNameActionFightPost(context.Background(), api.ActionFightMyNameActionFightPostParams{Name: c.name})
	if err != nil {
		return api.CharacterFightDataSchemaFight{}, err
	}

	switch v := res.(type) {
	case *api.CharacterFightResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		if v.Data.Fight.Result == api.CharacterFightDataSchemaFightResultLose {
			return api.CharacterFightDataSchemaFight{}, fmt.Errorf("loose battle")
		}

		return v.Data.Fight, c.updateData(unsafe.Pointer(&v.Data.Character))
	case *api.ActionFightMyNameActionFightPostCode486:
		return api.CharacterFightDataSchemaFight{}, fmt.Errorf("action is already in progress by your character")
	case *api.ActionFightMyNameActionFightPostCode497:
		return api.CharacterFightDataSchemaFight{}, fmt.Errorf("inventory is full")
	case *api.ActionFightMyNameActionFightPostCode498:
		return api.CharacterFightDataSchemaFight{}, fmt.Errorf("character not found")
	case *api.ActionFightMyNameActionFightPostCode499:
		return api.CharacterFightDataSchemaFight{}, fmt.Errorf("cooldown")
	case *api.ActionFightMyNameActionFightPostCode598:
		return api.CharacterFightDataSchemaFight{}, fmt.Errorf("monster is not at this map tile")
	default:
		return api.CharacterFightDataSchemaFight{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Craft(code string, quantity int) (api.SkillDataSchemaDetails, error) {
	res, err := c.cli.ActionCraftingMyNameActionCraftingPost(context.Background(), &api.CraftingSchema{Code: code, Quantity: api.NewOptInt(quantity)}, api.ActionCraftingMyNameActionCraftingPostParams{Name: c.name})
	if err != nil {
		return api.SkillDataSchemaDetails{}, err
	}

	switch v := res.(type) {
	case *api.SkillResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return v.Data.Details, c.updateData(unsafe.Pointer(&v.Data.Character))
	case *api.ActionCraftingMyNameActionCraftingPostNotFound:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("craft not found")
	case *api.ActionCraftingMyNameActionCraftingPostCode478:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("missing item or insufficient quantity")
	case *api.ActionCraftingMyNameActionCraftingPostCode486:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("action is already in progress by your character")
	case *api.ActionCraftingMyNameActionCraftingPostCode493:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("skill level is too low")
	case *api.ActionCraftingMyNameActionCraftingPostCode497:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("inventory is full")
	case *api.ActionCraftingMyNameActionCraftingPostCode498:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("character not found")
	case *api.ActionCraftingMyNameActionCraftingPostCode499:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("cooldown")
	case *api.ActionCraftingMyNameActionCraftingPostCode598:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("workshop not at this map tile")
	default:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Withdraw(code string, quantity int) error {
	res, err := c.cli.ActionWithdrawBankMyNameActionBankWithdrawPost(context.Background(), &api.SimpleItemSchema{Code: code, Quantity: quantity}, api.ActionWithdrawBankMyNameActionBankWithdrawPostParams{Name: c.name})
	if err != nil {
		return err
	}

	switch v := res.(type) {
	case *api.BankItemTransactionResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return c.updateData(unsafe.Pointer(&v.Data.Character))
	case *api.ActionWithdrawBankMyNameActionBankWithdrawPostNotFound:
		return fmt.Errorf("item not found")
	case *api.ActionWithdrawBankMyNameActionBankWithdrawPostCode461:
		return fmt.Errorf("transaction is already in progress with this item/your golds in your bank")
	case *api.ActionWithdrawBankMyNameActionBankWithdrawPostCode478:
		return fmt.Errorf("missing item or insufficient quantity")
	case *api.ActionWithdrawBankMyNameActionBankWithdrawPostCode486:
		return fmt.Errorf("action is already in progress by your character")
	case *api.ActionWithdrawBankMyNameActionBankWithdrawPostCode497:
		return fmt.Errorf("inventory is full")
	case *api.ActionWithdrawBankMyNameActionBankWithdrawPostCode498:
		return fmt.Errorf("character not found")
	case *api.ActionWithdrawBankMyNameActionBankWithdrawPostCode499:
		return fmt.Errorf("cooldown")
	case *api.ActionWithdrawBankMyNameActionBankWithdrawPostCode598:
		return fmt.Errorf("bank not at this map tile")
	default:
		return fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Deposit(code string, quantity int) error {
	res, err := c.cli.ActionDepositBankMyNameActionBankDepositPost(context.Background(), &api.SimpleItemSchema{Code: code, Quantity: quantity}, api.ActionDepositBankMyNameActionBankDepositPostParams{Name: c.name})
	if err != nil {
		return err
	}

	switch v := res.(type) {
	case *api.BankItemTransactionResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return c.updateData(unsafe.Pointer(&v.Data.Character))
	case *api.ActionDepositBankMyNameActionBankDepositPostNotFound:
		return fmt.Errorf("item not found")
	case *api.ActionDepositBankMyNameActionBankDepositPostCode461:
		return fmt.Errorf("transaction is already in progress with this item/your golds in your bank")
	case *api.ActionDepositBankMyNameActionBankDepositPostCode462:
		return fmt.Errorf("bank is full")
	case *api.ActionDepositBankMyNameActionBankDepositPostCode478:
		return fmt.Errorf("missing item or insufficient quantity")
	case *api.ActionDepositBankMyNameActionBankDepositPostCode486:
		return fmt.Errorf("action is already in progress by your character")
	case *api.ActionDepositBankMyNameActionBankDepositPostCode498:
		return fmt.Errorf("character not found")
	case *api.ActionDepositBankMyNameActionBankDepositPostCode499:
		return fmt.Errorf("cooldown")
	case *api.ActionDepositBankMyNameActionBankDepositPostCode598:
		return fmt.Errorf("bank not at this map tile")
	default:
		return fmt.Errorf("unknown answer type")
	}
}

func (c *Character) FindOnMap(code string) ([]api.MapSchema, error) {
	res, err := c.cli.GetAllMapsMapsGet(context.Background(), api.GetAllMapsMapsGetParams{ContentCode: api.NewOptString(code)})
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (c *Character) GetResource(code string) (api.ResourceSchema, error) {
	res, err := c.cli.GetResourceResourcesCodeGet(context.Background(), api.GetResourceResourcesCodeGetParams{Code: code})
	if err != nil {
		return api.ResourceSchema{}, err
	}

	switch v := res.(type) {
	case *api.ResourceResponseSchema:
		return v.Data, nil
	case *api.GetResourceResourcesCodeGetNotFound:
		return api.ResourceSchema{}, fmt.Errorf("resource not found")
	default:
		return api.ResourceSchema{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) GetMonster(code string) (api.MonsterSchema, error) {
	res, err := c.cli.GetMonsterMonstersCodeGet(context.Background(), api.GetMonsterMonstersCodeGetParams{Code: code})
	if err != nil {
		return api.MonsterSchema{}, err
	}

	switch v := res.(type) {
	case *api.MonsterResponseSchema:
		return v.Data, nil
	case *api.GetMonsterMonstersCodeGetNotFound:
		return api.MonsterSchema{}, fmt.Errorf("monster not found")
	default:
		return api.MonsterSchema{}, fmt.Errorf("unknown answer type: %v", v)
	}
}

func (c *Character) Recycle(code string, quantity int) (api.RecyclingDataSchemaDetails, error) {
	res, err := c.cli.ActionRecyclingMyNameActionRecyclingPost(context.Background(), &api.RecyclingSchema{Code: code, Quantity: api.NewOptInt(quantity)}, api.ActionRecyclingMyNameActionRecyclingPostParams{Name: c.name})
	if err != nil {
		return api.RecyclingDataSchemaDetails{}, err
	}

	switch v := res.(type) {
	case *api.RecyclingResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return v.Data.Details, c.updateData(unsafe.Pointer(&v.Data.Character))
	case *api.ActionRecyclingMyNameActionRecyclingPostNotFound:
		return api.RecyclingDataSchemaDetails{}, fmt.Errorf("item not found")
	case *api.ActionRecyclingMyNameActionRecyclingPostCode473:
		return api.RecyclingDataSchemaDetails{}, fmt.Errorf("item cannot be recycled")
	case *api.ActionRecyclingMyNameActionRecyclingPostCode478:
		return api.RecyclingDataSchemaDetails{}, fmt.Errorf("missing item or insufficient quantity")
	case *api.ActionRecyclingMyNameActionRecyclingPostCode486:
		return api.RecyclingDataSchemaDetails{}, fmt.Errorf("action is already in progress by your character")
	case *api.ActionRecyclingMyNameActionRecyclingPostCode493:
		return api.RecyclingDataSchemaDetails{}, fmt.Errorf("skill level is too low")
	case *api.ActionRecyclingMyNameActionRecyclingPostCode497:
		return api.RecyclingDataSchemaDetails{}, fmt.Errorf("inventory is full")
	case *api.ActionRecyclingMyNameActionRecyclingPostCode498:
		return api.RecyclingDataSchemaDetails{}, fmt.Errorf("character not found")
	case *api.ActionRecyclingMyNameActionRecyclingPostCode499:
		return api.RecyclingDataSchemaDetails{}, fmt.Errorf("cooldown")
	case *api.ActionRecyclingMyNameActionRecyclingPostCode598:
		return api.RecyclingDataSchemaDetails{}, fmt.Errorf("workshop not found on this map")
	default:
		return api.RecyclingDataSchemaDetails{}, fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Equip(code string, slot string, quantity int) error {
	res, err := c.cli.ActionEquipItemMyNameActionEquipPost(context.Background(), &api.EquipSchema{Code: code, Slot: api.EquipSchemaSlot(slot), Quantity: api.NewOptInt(quantity)}, api.ActionEquipItemMyNameActionEquipPostParams{Name: c.name})
	if err != nil {
		return err
	}

	switch v := res.(type) {
	case *api.EquipmentResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return c.updateData(unsafe.Pointer(&v.Data.Character))
	case *api.ActionEquipItemMyNameActionEquipPostNotFound:
		return fmt.Errorf("item not found")
	case *api.ActionEquipItemMyNameActionEquipPostCode478:
		return fmt.Errorf("missing item or insufficient quantity")
	case *api.ActionEquipItemMyNameActionEquipPostCode484:
		return fmt.Errorf("character can't equip more than 100 consumables in the same slot")
	case *api.ActionEquipItemMyNameActionEquipPostCode485:
		return fmt.Errorf("item is already equipped")
	case *api.ActionEquipItemMyNameActionEquipPostCode486:
		return fmt.Errorf("action is already in progress by your character")
	case *api.ActionEquipItemMyNameActionEquipPostCode491:
		return fmt.Errorf("slot is not empty")
	case *api.ActionEquipItemMyNameActionEquipPostCode496:
		return fmt.Errorf("character level is too low")
	case *api.ActionEquipItemMyNameActionEquipPostCode497:
		return fmt.Errorf("inventory is full")
	case *api.ActionEquipItemMyNameActionEquipPostCode498:
		return fmt.Errorf("character not found")
	case *api.ActionEquipItemMyNameActionEquipPostCode499:
		return fmt.Errorf("cooldown")
	default:
		return fmt.Errorf("unknown answer type: %v", v)
	}
}

func (c *Character) UnEquip(code string, slot string, quantity int) error {
	res, err := c.cli.ActionUnequipItemMyNameActionUnequipPost(context.Background(), &api.UnequipSchema{Slot: api.UnequipSchemaSlot(slot), Quantity: api.NewOptInt(quantity)}, api.ActionUnequipItemMyNameActionUnequipPostParams{Name: c.name})
	if err != nil {
		return err
	}

	switch v := res.(type) {
	case *api.EquipmentResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return c.updateData(unsafe.Pointer(&v.Data.Character))
	case *api.ActionUnequipItemMyNameActionUnequipPostNotFound:
		return fmt.Errorf("item not found")
	case *api.ActionUnequipItemMyNameActionUnequipPostCode478:
		return fmt.Errorf("missing item or insufficient quantity")
	case *api.ActionUnequipItemMyNameActionUnequipPostCode486:
		return fmt.Errorf("action is already in progress by your character")
	case *api.ActionUnequipItemMyNameActionUnequipPostCode491:
		return fmt.Errorf("slot is empty")
	case *api.ActionUnequipItemMyNameActionUnequipPostCode497:
		return fmt.Errorf("inventory is full")
	case *api.ActionUnequipItemMyNameActionUnequipPostCode498:
		return fmt.Errorf("character not found")
	case *api.ActionUnequipItemMyNameActionUnequipPostCode499:
		return fmt.Errorf("cooldown")
	default:
		return fmt.Errorf("unknown answer type: %v", v)
	}
}
