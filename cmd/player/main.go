package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Sinketsu/artifactsmmo/internal/characters/cetcalcoatl"
	"github.com/Sinketsu/artifactsmmo/internal/characters/enkidu"
	"github.com/Sinketsu/artifactsmmo/internal/characters/ereshkigal"
	"github.com/Sinketsu/artifactsmmo/internal/characters/ishtar"
	"github.com/Sinketsu/artifactsmmo/internal/events"
	"github.com/Sinketsu/artifactsmmo/internal/generic"
	"github.com/Sinketsu/artifactsmmo/internal/monitoring"
)

type Character interface {
	Live(ctx context.Context, events *events.Service)
}

func main() {
	go monitoring.NewClient(os.Getenv("MONITORING_WRITE_URL"), os.Getenv("MONITORING_FOLDER"), os.Getenv("MONITORING_TOKEN")).Run(30 * time.Second)

	serverParams := generic.ServerParams{
		ServerUrl:   os.Getenv("SERVER_URL"),
		ServerToken: os.Getenv("SERVER_TOKEN"),
	}

	events := events.New(serverParams)
	go events.Update(1 * time.Minute)

	characters := []Character{
		ishtar.NewCharacter(generic.Params{
			CharacterName: "Ishtar",
			ServerParams:  serverParams,
		}),
		cetcalcoatl.NewCharacter(generic.Params{
			CharacterName: "Cetcalcoatl",
			ServerParams:  serverParams,
		}),
		ereshkigal.NewCharacter(generic.Params{
			CharacterName: "Ereshkigal",
			ServerParams:  serverParams,
		}),
		enkidu.NewCharacter(generic.Params{
			CharacterName: "Enkidu",
			ServerParams:  serverParams,
		}),
	}

	ctx, stopNotify := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	wg := &sync.WaitGroup{}
	wg.Add(len(characters))
	for _, character := range characters {
		go func() {
			character.Live(ctx, events)
			wg.Done()
		}()
	}

	<-ctx.Done()
	fmt.Println("got stop signal...")

	stopNotify()
	wg.Wait()
}
