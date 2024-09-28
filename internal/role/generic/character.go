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
	name string
	data api.CharacterSchema

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
