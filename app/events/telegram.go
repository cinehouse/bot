package events

import (
	"context"
	"fmt"
	"github.com/cinehouse/bot/app/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

//go:generate mockery --inpackage --name=BotAPI --case=snake --testonly

// TelegramListener listens to tg update, forward to bots and send back responses
// Not thread safe
type TelegramListener struct {
	BotAPI BotAPI
	Bots   bot.Interface
	Debug  bool
}

type BotAPI interface {
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

// Do process all events, blocked call
func (l *TelegramListener) Do(ctx context.Context) (err error) {
	log.Printf("[TelegramListener] Start listening")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := l.BotAPI.GetUpdatesChan(u)

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
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			_, _ = l.BotAPI.Send(msg)
		}
	}
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
			ID: msg.From.ID,
		}
	}

	return &message
}
