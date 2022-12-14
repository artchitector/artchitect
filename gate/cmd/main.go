package main

import (
	"context"
	"github.com/artchitector/artchitect/gate/handler"
	"github.com/artchitector/artchitect/gate/repository"
	"github.com/artchitector/artchitect/gate/resources"
	"github.com/artchitector/artchitect/gate/state"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05"})

	res := resources.InitResources()
	log.Info().Msg("service gate started")

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		cancel()
	}()

	cardsRepository := repository.NewCardRepository(res.GetDB())
	decisionRepository := repository.NewDecisionRepository(res.GetDB())
	stateRepository := repository.NewStateRepository(res.GetDB())
	spellRepository := repository.NewSpellRepository(res.GetDB())
	lotteryRepository := repository.NewLotteryRepository(res.GetDB())
	retriever := state.NewRetriever(
		log.With().Str("service", "retriever").Logger(),
		cardsRepository,
		decisionRepository,
		stateRepository,
		spellRepository,
	)
	stateHandler := handler.NewStateHandler(
		log.With().Str("service", "state_handler").Logger(),
		retriever,
	)
	paintingHandler := handler.NewPaintingHandler(
		log.With().Str("service", "painting_handler").Logger(),
		retriever,
	)
	lastPaintingsHandler := handler.NewLastPaintingsHandler(cardsRepository)
	listHandler := handler.NewListHandler(cardsRepository)
	lotteryHandler := handler.NewLotteryHandler(
		log.With().Str("service", "lottery_handler").Logger(),
		lotteryRepository,
	)
	cardHandler := handler.NewCardHandler(cardsRepository)
	selectionHander := handler.NewSelectionHandler(lotteryRepository)

	go func() {
		r := gin.Default()
		r.Use(cors.New(cors.Config{
			AllowAllOrigins: true,
			//AllowOrigins:           nil,
		}))
		if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
			log.Fatal().Err(err).Send()
		}
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		r.GET("/state", stateHandler.Handle)
		r.GET("/painting/:id", paintingHandler.Handle)
		r.GET("/last_paintings/:quantity", lastPaintingsHandler.Handle)
		r.GET("/list/:from/:to", listHandler.Handle)
		r.GET("/lottery/:lastN", lotteryHandler.HandleLast)
		r.GET("/card/:id", cardHandler.Handle)
		r.GET("/selection", selectionHander.Handle)

		if err := r.Run("0.0.0.0:" + res.GetEnv().HttpPort); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	<-ctx.Done()
	log.Info().Msg("gate.Run finished")
}
