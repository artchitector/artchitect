package communication

import (
	"bytes"
	"context"
	"fmt"
	"gorm.io/gorm"
	"strings"

	"github.com/artchitector/artchitect/model"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type artPile interface {
	GetMaxArtID(ctx context.Context) (uint, error)
	GetArtRecursive(ctx context.Context, ID uint) (model.Art, error)
}

type settings interface {
	GetValue(ctx context.Context, name string) (string, error)
	SetValue(ctx context.Context, name string, value string) (model.Setting, error)
}

type warehouse interface {
	DownloadArtImage(ctx context.Context, artID uint, size string) ([]byte, error)
}

// Bot - телеграм-бот, который отправляет картинки в телеграм-чаты
type Bot struct {
	artPile                artPile
	warehouse              warehouse
	settings               settings
	bot                    *bot.Bot
	token                  string
	ChatArtchitectChoice   int64
	ChatArtchitectorChoice int64
	ArtchitectorID         int
}

func NewBot(
	artPile artPile,
	warehouse warehouse,
	settings settings,
	token string,
	chatArtchitectChoice int64,
	chatArtchitectorChoice int64,
	artchitectorID int,
) *Bot {
	b := &Bot{
		artPile:                artPile,
		warehouse:              warehouse,
		settings:               settings,
		token:                  token,
		ChatArtchitectChoice:   chatArtchitectChoice,
		ChatArtchitectorChoice: chatArtchitectorChoice,
		ArtchitectorID:         artchitectorID,
	}
	return b
}

func (b *Bot) Start(ctx context.Context) error {
	if b.token == "" {
		log.Info().Msgf("[bot] БОТ НЕ ИНИЦИАЛИЗИРУЕТСЯ!")
		return nil
	}
	log.Info().Msgf("[bot] НАЧАТА НАСТРОЙКА ТЕЛЕГРАМ-БОТА!")
	opts := []bot.Option{
		// bot.WithDebug(),
		// bot.WithCheckInitTimeout(time.Second * 10),
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
		if update.ChannelPost.ID == int(t.ArtchitectorID) {
			if err := t.handleArtchitector(ctx, update.Message.Text); err != nil {
				log.Error().Err(err).Msgf("[bot] НЕ СМОГ ОБРАБОТАТЬ ЗАПРОС ARTCHITEСTOR'А")
				return
			}
		}
	} else {
		log.Info().Msgf("[bot] ПОЛУЧЕНО ЛИЧНОЕ СООБЩЕНИЕ: %+v", update.Message)
	}
}

const (
	ArtchitectorSwitch = "переключи"
)

func (t *Bot) handleArtchitector(ctx context.Context, message string) error {
	if strings.Trim(message, " ") == ArtchitectorSwitch {
		// переключаем artchitect в ВЫКЛ/ВКЛ (запускается цикл творения или нет?)
		var newVal string
		val, err := t.settings.GetValue(ctx, model.SettingOdinActive)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newVal = model.OdinActive
		} else if err != nil {
			switch val {
			case model.OdinActive:
				newVal = model.OdinDisactive
			case model.OdinDisactive:
				newVal = model.OdinActive
			default:
				return fmt.Errorf("[bot] НЕИЗВЕСТНЫЙ СТАТУС НАСТРОЙКИ %s: %s", model.SettingOdinActive, val)
			}
		}
		if setting, err := t.settings.SetValue(ctx, model.SettingOdinActive, newVal); err != nil {
			return fmt.Errorf("[bot] ПРОБЛЕМЫ С СОХРАНЕНИЕМ НАСТРОЙКИ %s: %w", model.SettingOdinActive, err)
		} else {
			log.Info().Msgf("[bot] СОХРАНЕНА НАСТРОЙКА %s В ПОЛОЖЕНИИ %s", setting.SettingID, setting.Value)
			if err := t.replyArtchitector(ctx, fmt.Sprintf("[bot] СОХРАНЕНА НАСТРОЙКА %s В ПОЛОЖЕНИИ %s", setting.SettingID, setting.Value)); err != nil {
				return fmt.Errorf("[bot] НЕ МОГУ ОТПРАВИТЬ ОТВЕТ ARTCHITECTOR'У: %w", err)
			}
		}

	} else {
		if err := t.replyArtchitector(ctx, "НЕ ПОНИМАЮ КОМАНДУ"); err != nil {
			return fmt.Errorf("[bot] НЕ МОГУ ОТПРАВИТЬ ОТВЕТ ARTCHITECTOR'У: %w", err)
		}
	}
	return nil
}

func (b *Bot) replyArtchitector(ctx context.Context, message string) error {
	msg, err := b.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    b.ArtchitectorID,
		Text:      message,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		return fmt.Errorf("[bot] ПРОБЛЕМА ОТПРАВКИ СООБЩЕНИЯ ARTCHITECTOR'У: %w", err)
	}
	log.Info().Msgf("[bot] ОТПРАВИЛ СООБЩЕНИЕ ARTCHITECTOR'У: ID=%d", msg.ID)
	return nil
}
