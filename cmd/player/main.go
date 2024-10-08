package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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

func main() {
	cli := monitoring.NewClient(os.Getenv("MONITORING_WRITE_URL"), os.Getenv("MONITORING_FOLDER"), os.Getenv("MONITORING_TOKEN"))
	go cli.Run(30 * time.Second)

	serverParams := generic.ServerParams{
		ServerUrl:   os.Getenv("SERVER_URL"),
		ServerToken: os.Getenv("SERVER_TOKEN"),
	}

	events := events.New(serverParams)
	go events.Update(1 * time.Minute)

	Ishtar, err := ishtar.NewCharacter(generic.Params{
		CharacterName: "Ishtar",
		ServerParams:  serverParams,
	})
	if err != nil {
		panic(err)
	}

	Cetcalcoatl, err := cetcalcoatl.NewCharacter(generic.Params{
		CharacterName: "Cetcalcoatl",
		ServerParams:  serverParams,
	})
	if err != nil {
		panic(err)
	}

	Ereshkigal, err := ereshkigal.NewCharacter(generic.Params{
		CharacterName: "Ereshkigal",
		ServerParams:  serverParams,
	})
	if err != nil {
		panic(err)
	}

	Enkidu, err := enkidu.NewCharacter(generic.Params{
		CharacterName: "Enkidu",
		ServerParams:  serverParams,
	})
	if err != nil {
		panic(err)
	}

	ctx, stopNotify := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go Ishtar.Live(ctx, events)
	go Ereshkigal.Live(ctx, events)
	go Enkidu.Live(ctx, events)
	go Cetcalcoatl.Live(ctx, events)

	<-ctx.Done()
	fmt.Println("got stop signal...")

	stopNotify()
}
