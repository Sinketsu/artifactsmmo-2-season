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
	name string
	data api.CharacterSchema

	// cache FindOnMap
	mapsCache map[string]api.MapSchema
	// cache GetItem
	itemsCache map[string]api.SingleItemSchemaItem
	// cache GetMonster
	monsterCache map[string]api.MonsterSchema
	// cache GetResource
	resourceCache map[string]api.ResourceSchema

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

		mapsCache:     make(map[string]api.MapSchema),
		itemsCache:    make(map[string]api.SingleItemSchemaItem),
		monsterCache:  make(map[string]api.MonsterSchema),
		resourceCache: make(map[string]api.ResourceSchema),
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
	// tricky hack, because `ogen` generates different models for Character state from different methods instead of reusing one. But fields are the same - so we can cast it
	c.data = *(*api.CharacterSchema)(p)

	goldCount.Set(float64(c.data.Gold), c.name)
	tasksCoinCount.Set(float64(c.InInventory("tasks_coin")), c.name)
	return nil
}

func (c *Character) Data() api.CharacterSchema {
	return c.data
}
