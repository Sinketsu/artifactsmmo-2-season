package events

import (
	"context"
	"log/slog"
	"sync"
	"time"

	oas "github.com/Sinketsu/artifactsmmo/gen/oas"
	"github.com/Sinketsu/artifactsmmo/internal/api"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

type Service struct {
	cli    *api.Client
	logger *slog.Logger

	events []oas.ActiveEventSchema
	mu     sync.Mutex
}

func New(client *api.Client) *Service {
	return &Service{
		cli:    client,
		logger: slog.Default().With(ycloggingslog.Stream, "events"),
	}
}

func (s *Service) Update(ctx context.Context, interval time.Duration) {
	if err := s.update(); err != nil {
		s.logger.With(slog.Any("error", err)).Error("fail update event list")
		errorRate.Inc()
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.update(); err != nil {
				s.logger.With(slog.Any("error", err)).Error("fail update event list")
				errorRate.Inc()
			}
		}
	}
}

func (s *Service) update() error {
	result, err := s.cli.GetAllEventsEventsGet(context.TODO(), oas.GetAllEventsEventsGetParams{})
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.events = result.Data
	s.mu.Unlock()

	eventCount.Set(int64(len(s.events)))
	return nil
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
