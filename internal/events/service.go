package events

import (
	"context"
	"fmt"
	"sync"
	"time"

	oas "github.com/Sinketsu/artifactsmmo/gen/oas"
	"github.com/Sinketsu/artifactsmmo/internal/api"
)

type Service struct {
	cli *api.Client

	events []oas.ActiveEventSchema
	mu     sync.Mutex
}

func New(client *api.Client) *Service {
	return &Service{
		cli: client,
	}
}

func (s *Service) Update(interval time.Duration) {
	for range time.Tick(interval) {
		result, err := s.cli.GetAllEventsEventsGet(context.TODO(), oas.GetAllEventsEventsGetParams{})
		if err != nil {
			fmt.Println("fail update event list:", err)
		}

		s.mu.Lock()
		s.events = result.Data
		s.mu.Unlock()
	}
}

func (s *Service) Get(name string) *oas.ActiveEventSchema {
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
