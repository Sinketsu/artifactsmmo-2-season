package events

import (
	"context"
	"fmt"
	"sync"
	"time"

	api "github.com/Sinketsu/artifactsmmo/gen/oas"
	"github.com/Sinketsu/artifactsmmo/internal/generic"
)

type Service struct {
	cli *api.Client

	events []api.ActiveEventSchema
	mu     sync.Mutex
}

func New(params generic.ServerParams) *Service {
	client, err := api.NewClient(params.ServerUrl, &generic.Auth{Token: params.ServerToken})
	if err != nil {
		panic(err)
	}

	return &Service{
		cli: client,
	}
}

func (s *Service) Update(interval time.Duration) {
	for range time.Tick(interval) {
		result, err := s.cli.GetAllEventsEventsGet(context.TODO(), api.GetAllEventsEventsGetParams{})
		if err != nil {
			fmt.Println("fail update event list:", err)
		}

		s.mu.Lock()
		s.events = result.Data
		s.mu.Unlock()
	}
}

func (s *Service) Get(name string) *api.ActiveEventSchema {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ev := range s.events {
		if ev.Name == name {
			if ev.Expiration.After(time.Now()) {
				return &ev
			}
		}
	}

	return nil
}
