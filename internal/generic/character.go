package generic

import (
	"context"
	"fmt"
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
		return c.updateData(unsafe.Pointer(&v.Data))
	case *api.GetCharacterCharactersNameGetNotFound:
		return fmt.Errorf("character not found")
	default:
		return fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Log(msg ...any) {
	fmt.Print(c.name + ": ")
	fmt.Println(msg...)
}

func (c *Character) updateData(p unsafe.Pointer) error {
	c.data = *(*api.CharacterSchema)(p)

	goldCount.Set(float64(c.data.Gold), c.name)
	tasksCoinCount.Set(float64(c.InInventory("tasks_coin")), c.name)
	return nil
}

func (c *Character) Data() api.CharacterSchema {
	return c.data
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
