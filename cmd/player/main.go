package main

import (
	"context"
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
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
	ycmonitoringgo "github.com/Sinketsu/yc-monitoring-go"
	ycsdk "github.com/yandex-cloud/go-sdk"
)

type Character interface {
	Live(ctx context.Context)
}

func main() {
	logHandler, err := ycloggingslog.New(ycloggingslog.Options{
		LogGroupId:   os.Getenv("LOGGING_GROUP_ID"),
		ResourceType: "app",
		ResourceId:   "season-2",
		Credentials:  ycsdk.OAuthToken(os.Getenv("LOGGING_TOKEN")),
	})
	if err != nil {
		slog.With(slog.Any("error", err)).Error("fail to init log handler")
		os.Exit(1)
	}
	slog.SetDefault(slog.New(logHandler))

	apiClient, err := api.NewClient(api.Params{
		ServerUrl:   os.Getenv("SERVER_URL"),
		ServerToken: os.Getenv("SERVER_TOKEN"),
	})
	if err != nil {
		slog.With(slog.Any("error", err)).Error("fail to init API client")
		os.Exit(1)
	}

	bank := bank.New(apiClient)
	events := events.New(apiClient)
	monitoringClient := ycmonitoringgo.NewClient(os.Getenv("MONITORING_FOLDER"), os.Getenv("MONITORING_TOKEN"), ycmonitoringgo.WithLogger(slog.Default()))

	characters := []Character{
		ishtar.NewCharacter(apiClient, bank, events),
		cetcalcoatl.NewCharacter(apiClient, bank, events),
		ereshkigal.NewCharacter(apiClient, bank, events),
		enkidu.NewCharacter(apiClient, bank, events),
	}

	ctx, stopNotify := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go events.Update(ctx, 1*time.Minute)
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
	slog.Info("got stop signal...")

	stopNotify()
	wg.Wait()
}
