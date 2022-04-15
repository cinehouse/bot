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
	messageLogger := &MockMessageLogger{}
	tgBotAPI := &MockTgBotAPI{}
	bots := &bot.MockInterface{}

	l := TelegramListener{
		TgBotAPI:      tgBotAPI,
		MessageLogger: messageLogger,
		Bots:          bots,
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

	tgBotAPI.On("GetChat", mock.Anything).Return(tgbotapi.Chat{ID: 123}, nil)

	updChan := make(chan tgbotapi.Update, 1)
	updChan <- updMsg
	close(updChan)
	tgBotAPI.On("GetUpdatesChan", mock.Anything).Return(tgbotapi.UpdatesChannel(updChan), nil)

	bots.On("OnMessage", mock.Anything).Return(bot.Response{Send: false})
	messageLogger.On("Save", mock.MatchedBy(func(msg *bot.Message) bool {
		t.Logf("%v", msg)
		return msg.Text == "text 123" && msg.From.UserName == "user"
	}))
	err := l.Do(ctx)
	assert.EqualError(t, err, "[ERROR] Telegram listener: updates channel closed")

	messageLogger.AssertExpectations(t)
	messageLogger.AssertNumberOfCalls(t, "Save", 1)
}

func TestTelegramListener_DoWithBots(t *testing.T) {
	messageLogger := &MockMessageLogger{}
	tgBotAPI := &MockTgBotAPI{}
	bots := &bot.MockInterface{}

	l := TelegramListener{
		MessageLogger: messageLogger,
		TgBotAPI:      tgBotAPI,
		Bots:          bots,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Minute)
	defer cancel()

	updMsg := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: 123},
			Text: "text 123",
			From: &tgbotapi.User{UserName: "user"},
			Date: int(time.Date(2020, 2, 11, 19, 35, 55, 9, time.UTC).Unix()),
		},
	}

	tgBotAPI.On("GetChat", mock.Anything).Return(tgbotapi.Chat{ID: 123}, nil)

	updChan := make(chan tgbotapi.Update, 1)
	updChan <- updMsg
	close(updChan)
	tgBotAPI.On("GetUpdatesChan", mock.Anything).Return(tgbotapi.UpdatesChannel(updChan), nil)

	bots.On("OnMessage", mock.MatchedBy(func(msg bot.Message) bool {
		t.Logf("on-message: %+v", msg)
		return msg.Text == "text 123" && msg.From.UserName == "user"
	})).Return(bot.Response{Send: true, Text: "bot's answer"})

	messageLogger.On("Save", mock.MatchedBy(func(msg *bot.Message) bool {
		t.Logf("save: %+v", msg)
		return msg.Text == "text 123" && msg.From.UserName == "user"
	}))
	messageLogger.On("Save", mock.MatchedBy(func(msg *bot.Message) bool {
		t.Logf("save: %+v", msg)
		return msg.Text == "bot's answer"
	}))

	tgBotAPI.On("Send", mock.MatchedBy(func(c tgbotapi.MessageConfig) bool {
		t.Logf("send: %+v", c)
		return c.Text == "bot's answer"
	})).Return(tgbotapi.Message{Text: "bot's answer", From: &tgbotapi.User{UserName: "user"}}, nil)

	err := l.Do(ctx)
	assert.EqualError(t, err, "[ERROR] Telegram listener: updates channel closed")
	messageLogger.AssertExpectations(t)
	messageLogger.AssertNumberOfCalls(t, "Save", 2)
	tgBotAPI.AssertNumberOfCalls(t, "Send", 1)
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
