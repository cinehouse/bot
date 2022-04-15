package events

import (
	"context"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cinehouse/bot/app/bot"
)

func TestTelegramListener_DoNoBots(t *testing.T) {
	tbAPI := &MockBotAPI{}
	bots := &bot.MockInterface{}

	l := TelegramListener{
		BotAPI: tbAPI,
		Bots:   bots,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	updMsg := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: 123},
			Text: "text 123",
			From: &tgbotapi.User{UserName: "user"},
		},
	}

	tbAPI.On("GetChat", mock.Anything).Return(tgbotapi.Chat{ID: 123}, nil)

	updChan := make(chan tgbotapi.Update, 1)
	updChan <- updMsg
	close(updChan)
	tbAPI.On("GetUpdatesChan", mock.Anything).Return(tgbotapi.UpdatesChannel(updChan), nil)

	bots.On("OnMessage", mock.Anything).Return(bot.Response{Send: false})
	err := l.Do(ctx)
	assert.EqualError(t, err, "[ERROR] Telegram listener: updates channel closed")
}

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
