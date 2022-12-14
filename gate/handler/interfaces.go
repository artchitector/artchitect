package handler

import (
	"context"
	"github.com/artchitector/artchitect/model"
)

type retriever interface {
	CollectState(ctx context.Context) (model.CurrentState, error)
	GetPaintingData(ctx context.Context, paintingID uint) ([]byte, error)
}

type cardsRepository interface {
	GetCard(ctx context.Context, ID uint) (model.Card, bool, error)
	GetLastCards(ctx context.Context, count uint64) ([]model.Card, error)
	GetCardsRange(ctx context.Context, from uint, to uint) ([]model.Card, error)
}

type lotteryRepository interface {
	GetActiveLottery(ctx context.Context) (model.Lottery, error)
	GetLastLotteries(ctx context.Context, lastN uint) ([]model.Lottery, error)
	GetSelection(ctx context.Context) ([]uint64, error)
}
