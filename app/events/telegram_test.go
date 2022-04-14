package events

import (
	"github.com/cinehouse/bot/app/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTelegram_transformTextMessage(t *testing.T) {
	l := TelegramListener{}
	assert.Equal(t,
		&bot.Message{
			MessageID: 30,
			From: bot.User{
				ID: 123,
			},
			Chat: bot.Chat{
				ID: 456,
			},
			Text: "test",
		},
		l.transform(
			&tgbotapi.Message{
				MessageID: 30,
				From: &tgbotapi.User{
					ID: 123,
				},
				Chat: &tgbotapi.Chat{
					ID: 456,
				},
				Text: "test",
			},
		),
	)
}
