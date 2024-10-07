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

	go Ishtar.Live(ctx)
	go Ereshkigal.Live(ctx)
	go Enkidu.Live(ctx)
	go Cetcalcoatl.Live(ctx)

	<-ctx.Done()
	fmt.Println("got stop signal...")

	stopNotify()
}
