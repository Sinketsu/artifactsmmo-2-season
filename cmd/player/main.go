package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sinketsu/artifactsmmo/internal/role/fighter"
	"github.com/Sinketsu/artifactsmmo/internal/role/gatherer"
	"github.com/Sinketsu/artifactsmmo/internal/role/generic"
)

func main() {
	serverParams := generic.ServerParams{
		ServerUrl:   os.Getenv("SERVER_URL"),
		ServerToken: os.Getenv("SERVER_TOKEN"),
	}

	ishtar, err := gatherer.NewCharacter(generic.Params{
		CharacterName: "Ishtar",
		ServerParams:  serverParams,
	})
	if err != nil {
		panic(err)
	}

	ereshkigal, err := fighter.NewCharacter(generic.Params{
		CharacterName: "Ereshkigal",
		ServerParams:  serverParams,
	})
	if err != nil {
		panic(err)
	}

	ctx, stopNotify := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go ishtar.Live(ctx)
	go ereshkigal.Live(ctx)

	<-ctx.Done()
	fmt.Println("got stop signal...")

	stopNotify()
}
