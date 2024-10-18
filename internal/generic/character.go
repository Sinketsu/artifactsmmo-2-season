package generic

import (
	"context"
	"fmt"
	"time"
	"unsafe"

	oas "github.com/Sinketsu/artifactsmmo/gen/oas"
	"github.com/Sinketsu/artifactsmmo/internal/api"
	"github.com/Sinketsu/artifactsmmo/internal/bank"
)

// Bank is not thread safe - so you need to explicit call Lock() and Unlock()
type Bank interface {
	Lock()
	Unlock()
	Items() ([]bank.Item, error)
}

type Events interface {
	Get(name string) *oas.ActiveEventSchema
}

type Params struct {
	Name string
}

type Character struct {
	name string
	data oas.CharacterSchema

	bank   Bank
	events Events

	// cache FindOnMap
	mapsCache map[string]oas.MapSchema
	// cache GetItem
	itemsCache map[string]oas.SingleItemSchemaItem
	// cache GetMonster
	monsterCache map[string]oas.MonsterSchema
	// cache GetResource
	resourceCache map[string]oas.ResourceSchema

	cli *api.Client
}

func NewCharacter(client *api.Client, params Params, bank Bank, events Events) (*Character, error) {
	character := &Character{
		name: params.Name,
		cli:  client,

		bank:   bank,
		events: events,

		mapsCache:     make(map[string]oas.MapSchema),
		itemsCache:    make(map[string]oas.SingleItemSchemaItem),
		monsterCache:  make(map[string]oas.MonsterSchema),
		resourceCache: make(map[string]oas.ResourceSchema),
	}

	if err := character.initData(); err != nil {
		return nil, err
	}

	return character, nil
}

func (c *Character) Name() string {
	return c.name
}

func (c *Character) Bank() Bank {
	return c.bank
}

func (c *Character) Events() Events {
	return c.events
}

func (c *Character) initData() error {
	res, err := c.cli.GetCharacterCharactersNameGet(context.Background(), oas.GetCharacterCharactersNameGetParams{Name: c.name})
	if err != nil {
		return err
	}

	switch v := res.(type) {
	case *oas.CharacterResponseSchema:
		return c.updateData(unsafe.Pointer(&v.Data))
	case *oas.GetCharacterCharactersNameGetNotFound:
		return fmt.Errorf("character not found")
	default:
		return fmt.Errorf("unknown answer type")
	}
}

func (c *Character) Log(msg ...any) {
	fmt.Print("["+time.Now().Local().Format(time.TimeOnly)+"] ", c.name+": ")
	fmt.Println(msg...)
}

func (c *Character) updateData(p unsafe.Pointer) error {
	// tricky hack, because `ogen` generates different models for Character state from different methods instead of reusing one. But fields are the same - so we can cast it
	c.data = *(*oas.CharacterSchema)(p)

	goldCount.Set(float64(c.data.Gold), c.name)
	tasksCoinCount.Set(float64(c.InInventory("tasks_coin")), c.name)

	skillLevel.Set(float64(c.data.GearcraftingLevel), c.name, "gearcrafting")
	skillLevel.Set(float64(c.data.WeaponcraftingLevel), c.name, "weaponcrafting")
	skillLevel.Set(float64(c.data.JewelrycraftingLevel), c.name, "jewerlycrafting")
	skillLevel.Set(float64(c.data.CookingLevel), c.name, "cooking")
	skillLevel.Set(float64(c.data.MiningLevel), c.name, "mining")
	skillLevel.Set(float64(c.data.WoodcuttingLevel), c.name, "woodcutting")
	skillLevel.Set(float64(c.data.FishingLevel), c.name, "fishing")

	return nil
}

func (c *Character) Data() oas.CharacterSchema {
	return c.data
}
