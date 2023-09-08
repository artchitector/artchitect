package communication

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artchitector/artchitect2/model"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"
)

// Bot - телеграм-бот, который отправляет картинки в телеграм-чаты
type Bot struct {
	artPile                artPile
	warehouse              warehouse
	bot                    *bot.Bot
	token                  string
	ChatArtchitectChoice   int64
	ChatArtchitectorChoice int64
}

func NewBot(
	artPile artPile,
	warehouse warehouse,
	token string,
	chatArtchitectChoice int64,
	chatArtchitectorChoice int64,
) *Bot {
	b := &Bot{
		artPile:                artPile,
		warehouse:              warehouse,
		token:                  token,
		ChatArtchitectChoice:   chatArtchitectChoice,
		ChatArtchitectorChoice: chatArtchitectorChoice,
	}
	return b
}

func (b *Bot) Start(ctx context.Context) error {
	log.Info().Msgf("[bot] НАЧАТА НАСТРОЙКА ТЕЛЕГРАМ-БОТА!")
	opts := []bot.Option{
		//bot.WithDebug(),
		//bot.WithCheckInitTimeout(time.Second * 10),
		bot.WithDefaultHandler(b.defaultHandler),
	}
	if tgBot, err := bot.New(b.token, opts...); err != nil {
		return errors.Wrapf(err, "[bot] НЕ МОГУ ПОДКЛЮЧИТЬСЯ К ТЕЛЕГРАМ-БОТУ %s!", b.token)
	} else {
		// start bot to listen all messages
		log.Info().Msgf("[bot] НАСТРОЙКА БОТА ПРОИЗВЕДЕНА!")
		b.bot = tgBot
	}

	go func() {
		log.Info().Msgf("[bot] НАЧИНАЮ ПРОСЛУШИВАНИЕ ТЕЛЕГРАМ-БОТА")
		b.bot.Start(ctx)
		log.Info().Msgf("[bot] БОТ ЗАВЕРШЁН")
	}()

	return nil
}

func (b *Bot) SendArtchitectChoice(ctx context.Context, artID uint) error {
	return b.send2Chat(ctx, artID, b.ChatArtchitectChoice)
}

func (b *Bot) SendArtchitectorChoice(ctx context.Context, artID uint) error {
	return b.send2Chat(ctx, artID, b.ChatArtchitectorChoice)
}

func (b *Bot) send2Chat(ctx context.Context, artID uint, chatID int64) error {
	if b.bot == nil {
		return errors.Errorf("[bot] БОТ: НЕ ИНИЦИАЛИЗИРОВАНО")
	}

	text, err := b.generateText(ctx, artID)
	if err != nil {
		return errors.Wrapf(err, "[bot] НЕ МОГУ СГЕНЕРИРОВАТЬ ТЕКСТ ДЛЯ #%d", artID)
	}

	img, err := b.downloadImage(ctx, artID)
	if err != nil {
		return errors.Wrapf(err, "[bot] НЕ МОГУ ЗАГРУЗИТЬ КАРТИНКУ #%d", artID)
	}

	r := bytes.NewReader(img)

	log.Info().Msgf("[bot] ОТПРАВЛЯЮ СООБЩЕНИЕ В ЧАТ %d: %s", chatID, strings.ReplaceAll(text, "\n", "\\n"))
	msg, err := b.bot.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:    chatID,
		Photo:     &models.InputFileUpload{Data: r},
		Caption:   text,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		return errors.Wrapf(err, "[bot] НЕ МОГУ ОТПРАВИТЬ СООБЩЕНИЕ С КАРТИНОЙ #%d", artID)
	}

	log.Info().Msgf("[bot] СООБЩЕНИЕ ОТПРАВЛЕНО. УСПЕХ. ID=%d", msg.ID)
	return nil
}

func (b *Bot) downloadImage(ctx context.Context, artID uint) ([]byte, error) {
	return b.warehouse.DownloadArtImage(ctx, artID, model.SizeF)
}

func (b *Bot) generateText(ctx context.Context, artID uint) (string, error) {
	art, err := b.artPile.GetArtRecursive(ctx, artID)
	if err != nil {
		return "", errors.Wrapf(err, "[bot] НЕ МОГУ ПОЛУЧИТЬ ART #%d", artID)
	}

	return fmt.Sprintf(
		"<b>Art <a href='https://artchitect.space/art/%d'>#%d</a></b>.\n"+
			"created: <i>%s</i>, seed: <b>%d</b>\n"+
			"words: <i>%s</i>",
		art.ID,
		art.ID,
		art.CreatedAt.Format("2006 Jan 2 15:04"),
		art.Idea.Seed,
		art.Idea.WordsStr,
	), nil
}

func (t *Bot) defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.ChannelPost != nil {
		log.Info().Msgf("[bot] ПОЛУЧЕНО СООБЩЕНИЕ В КАНАЛЕ: %+v", update.ChannelPost)
	} else {
		log.Info().Msgf("[bot] ПОЛУЧЕНО ЛИЧНОЕ СООБЩЕНИЕ: %+v", update.Message)
	}
}
