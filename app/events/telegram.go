package events

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"

	"github.com/cinehouse/bot/app/bot"
)

//go:generate mockery --inpackage --name=TgBotAPI --case=snake
//go:generate mockery --inpackage --name=MessageLogger --case=snake

// TelegramListener listens to tg update, forward to bots and send back responses
// Not thread safe
type TelegramListener struct {
	TgBotAPI      TgBotAPI
	MessageLogger MessageLogger
	Bots          bot.Interface
	SuperUsers    SuperUser
	Debug         bool
}

type TgBotAPI interface {
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type MessageLogger interface {
	Save(msg *bot.Message)
}

// Do process all events, blocked call
func (l *TelegramListener) Do(ctx context.Context) (err error) {
	log.Printf("[TelegramListener] Start listening")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := l.TgBotAPI.GetUpdatesChan(u)

	for {
		select {

		case <-ctx.Done():
			return ctx.Err()

		case update, ok := <-updates:
			if !ok {
				return fmt.Errorf("[ERROR] Telegram listener: updates channel closed")
			}

			if l.Debug {
				log.Printf("[DEBUG] Update: %+v", update)
			}

			if update.Message == nil {
				log.Printf("[ERROR] Telegram listener: update.Message is nil")
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := l.transform(update.Message)
			l.MessageLogger.Save(msg) // save an incoming update to report

			log.Printf("[DEBUG] Incoming msg: %+v", msg)

			resp := l.Bots.OnMessage(*msg)

			if err := l.sendBotResponse(update.Message.Chat.ID, resp); err != nil {
				log.Printf("[WARN] failed to respond on update, %v", err)
			}
		}
	}
}

func (l *TelegramListener) sendBotResponse(chatID int64, resp bot.Response) error {
	if !resp.Send {
		return nil
	}

	log.Printf("[DEBUG] Bot response - %+v", resp.Text)

	msg := tgbotapi.NewMessage(chatID, resp.Text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	msg.DisableWebPagePreview = !resp.Preview

	res, err := l.TgBotAPI.Send(msg)
	if err != nil {
		return errors.Wrapf(err, "can't send message to telegram %q", resp.Text)
	}

	l.saveBotMessage(&res)

	return nil
}

func (l *TelegramListener) saveBotMessage(msg *tgbotapi.Message) {
	l.MessageLogger.Save(l.transform(msg))
}

func (l *TelegramListener) transform(msg *tgbotapi.Message) *bot.Message {
	message := bot.Message{
		MessageID: msg.MessageID,
		Text:      msg.Text,
	}

	if msg.Chat != nil {
		message.Chat = bot.Chat{
			ID: msg.Chat.ID,
		}
	}

	if msg.From != nil {
		message.From = bot.User{
			ID:       msg.From.ID,
			UserName: msg.From.UserName,
		}
	}

	return &message
}
