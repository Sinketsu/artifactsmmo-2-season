package generic

import (
	"context"
	"fmt"
	"time"
	"unsafe"

	api "github.com/Sinketsu/artifactsmmo/gen/oas"
)

type ServerParams struct {
	ServerUrl   string
	ServerToken string
}

type Params struct {
	ServerParams ServerParams

	CharacterName string
}

type Character struct {
	name       string
	data       api.CharacterSchema
	gatherData GaterData
	fightData  FightData
	craftData  CraftData

	cli *api.Client
}

func NewCharacter(params Params) (*Character, error) {
	client, err := api.NewClient(params.ServerParams.ServerUrl, &Auth{Token: params.ServerParams.ServerToken})
	if err != nil {
		return nil, err
	}

	character := &Character{
		name: params.CharacterName,
		cli:  client,
	}

	if err := character.initData(); err != nil {
		return nil, err
	}

	return character, nil
}

func (c *Character) Name() string {
	return c.name
}

func (c *Character) initData() error {
	res, err := c.cli.GetCharacterCharactersNameGet(context.Background(), api.GetCharacterCharactersNameGetParams{Name: c.name})
	if err != nil {
		return err
	}

	switch v := res.(type) {
	case *api.CharacterResponseSchema:
		c.data = v.Data
		return nil
	default:
		return fmt.Errorf("unknown answer type: %v", v)
	}
}

func (c *Character) Log(msg ...any) {
	fmt.Print(c.name + ": ")
	fmt.Println(msg...)
}

func (c *Character) updateData(p unsafe.Pointer) error {
	c.data = *(*api.CharacterSchema)(p)

	return nil
}

func (c *Character) Data() api.CharacterSchema {
	return c.data
}

func (c *Character) Gather() (api.SkillDataSchemaDetails, error) {
	res, err := c.cli.ActionGatheringMyNameActionGatheringPost(context.Background(), api.ActionGatheringMyNameActionGatheringPostParams{Name: c.name})
	if err != nil {
		return api.SkillDataSchemaDetails{}, err
	}

	switch v := res.(type) {
	case *api.SkillResponseSchema:
		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)

		return v.Data.Details, c.updateData(unsafe.Pointer(&v.Data.Character))
	default:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("unknown answer type: %v", v)
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
	default:
		return fmt.Errorf("unknown answer type: %v", v)
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
	default:
		return 0, fmt.Errorf("unknown answer type: %v", v)
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
	default:
		return api.GEItemSchema{}, fmt.Errorf("unknown answer type: %v", v)
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
	default:
		return api.SingleItemSchemaItem{}, fmt.Errorf("unknown answer type: %v", v)
	}
}

func (c *Character) InventoryItemCount() int {
	count := 0
	for _, item := range c.data.Inventory {
		count += item.Quantity
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
	default:
		return api.CharacterFightDataSchemaFight{}, fmt.Errorf("unknown answer type: %v", v)
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
	default:
		return api.SkillDataSchemaDetails{}, fmt.Errorf("unknown answer type: %v", v)
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
	default:
		return fmt.Errorf("unknown answer type: %v", v)
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
	default:
		return fmt.Errorf("unknown answer type: %v", v)
	}
}

func (c *Character) FindOnMap(code string) ([]api.MapSchema, error) {
	res, err := c.cli.GetAllMapsMapsGet(context.Background(), api.GetAllMapsMapsGetParams{ContentCode: api.NewOptString(code)})
	if err != nil {
		return nil, err
	}

	return res.Data, nil
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
	default:
		return api.RecyclingDataSchemaDetails{}, fmt.Errorf("unknown answer type: %v", v)
	}
}
