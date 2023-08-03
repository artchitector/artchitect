package main

import (
	"context"
	"github.com/artchitector/artchitect2/services/soul/infrastructure"
	"github.com/artchitector/artchitect2/services/soul/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05"})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		cancel()
	}()

	config := infrastructure.InitEnv()

	go runServices(ctx)
	runArtchitect(ctx, config)
}

// runArtchitect - main creative loop
func runArtchitect(ctx context.Context, config *infrastructure.Config) {
	creator := internal.NewCreator(config.CreatorActive)
	artchitect := internal.NewArtchitect(creator)
	artchitect.Run(ctx)
}

// runServices - periodically services (run once in minute etc.)
func runServices(ctx context.Context) {

}
