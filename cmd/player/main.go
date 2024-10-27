package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Sinketsu/artifactsmmo/internal/api"
	"github.com/Sinketsu/artifactsmmo/internal/bank"
	"github.com/Sinketsu/artifactsmmo/internal/characters/cetcalcoatl"
	"github.com/Sinketsu/artifactsmmo/internal/characters/enkidu"
	"github.com/Sinketsu/artifactsmmo/internal/characters/ereshkigal"
	"github.com/Sinketsu/artifactsmmo/internal/characters/ishtar"
	"github.com/Sinketsu/artifactsmmo/internal/events"
	ycmonitoringgo "github.com/Sinketsu/yc-monitoring-go"
)

type Character interface {
	Live(ctx context.Context)
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	monitoringClient := ycmonitoringgo.NewClient(os.Getenv("MONITORING_FOLDER"), os.Getenv("MONITORING_TOKEN"), ycmonitoringgo.WithLogger(logger))

	apiClient, err := api.NewClient(api.Params{
		ServerUrl:   os.Getenv("SERVER_URL"),
		ServerToken: os.Getenv("SERVER_TOKEN"),
	})
	if err != nil {
		panic(err)
	}

	bank := bank.New(apiClient)
	events := events.New(apiClient)
	go events.Update(1 * time.Minute)

	characters := []Character{
		ishtar.NewCharacter(apiClient, bank, events, os.Getenv("LOGGING_GROUP_ID"), os.Getenv("LOGGING_TOKEN")),
		cetcalcoatl.NewCharacter(apiClient, bank, events, os.Getenv("LOGGING_GROUP_ID"), os.Getenv("LOGGING_TOKEN")),
		ereshkigal.NewCharacter(apiClient, bank, events, os.Getenv("LOGGING_GROUP_ID"), os.Getenv("LOGGING_TOKEN")),
		enkidu.NewCharacter(apiClient, bank, events, os.Getenv("LOGGING_GROUP_ID"), os.Getenv("LOGGING_TOKEN")),
	}

	ctx, stopNotify := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go monitoringClient.Run(ctx, ycmonitoringgo.DefaultRegistry, 30*time.Second)

	wg := &sync.WaitGroup{}
	wg.Add(len(characters))
	for _, character := range characters {
		go func() {
			character.Live(ctx)
			wg.Done()
		}()
	}

	<-ctx.Done()
	fmt.Println("got stop signal...")

	stopNotify()
	wg.Wait()
}
