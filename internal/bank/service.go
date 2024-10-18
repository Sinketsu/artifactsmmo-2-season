package bank

import (
	"context"
	"fmt"
	"sync"

	oas "github.com/Sinketsu/artifactsmmo/gen/oas"
	"github.com/Sinketsu/artifactsmmo/internal/api"
)

type Item struct {
	oas.SingleItemSchemaItem
	Quantity int
}

type Service struct {
	cli *api.Client

	mu sync.Mutex

	// more caches to god of caches!
	itemCache map[string]oas.SingleItemSchemaItem
}

func New(client *api.Client) *Service {
	return &Service{
		cli: client,

		itemCache: make(map[string]oas.SingleItemSchemaItem),
	}
}

func (s *Service) Lock() {
	s.mu.Lock()
}

func (s *Service) Unlock() {
	s.mu.Unlock()
}

func (s *Service) Items() ([]Item, error) {
	simpleResult, err := s.cli.GetBankItemsMyBankItemsGet(context.Background(), oas.GetBankItemsMyBankItemsGetParams{Size: oas.NewOptInt(100)})
	if err != nil {
		return nil, err
	}

	result := make([]Item, len(simpleResult.Data))
	for _, simpleItem := range simpleResult.Data {
		item, err := s.getItem(simpleItem.Code)
		if err != nil {
			return nil, err
		}

		result = append(result, Item{SingleItemSchemaItem: item, Quantity: simpleItem.Quantity})
	}

	return result, nil
}

func (s *Service) getItem(code string) (oas.SingleItemSchemaItem, error) {
	if item, ok := s.itemCache[code]; ok {
		return item, nil
	}

	res, err := s.cli.GetItemItemsCodeGet(context.Background(), oas.GetItemItemsCodeGetParams{Code: code})
	if err != nil {
		return oas.SingleItemSchemaItem{}, err
	}

	switch v := res.(type) {
	case *oas.ItemResponseSchema:
		s.itemCache[code] = v.Data.Item
		return v.Data.Item, nil
	case *oas.GetItemItemsCodeGetNotFound:
		return oas.SingleItemSchemaItem{}, fmt.Errorf("item not found")
	default:
		return oas.SingleItemSchemaItem{}, fmt.Errorf("unknown answer type")
	}
}
