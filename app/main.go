package main

import (
	"context"
	"fmt"
	"github.com/cinehouse/bot/app/events"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

var opts struct {
	Telegram struct {
		Token string `long:"token" env:"TOKEN" description:"Telegram bot token"`
	} `group:"telegram" namespace:"telegram" env-namespace:"TELEGRAM"`
	Debug bool `long:"debug" description:"Show debug information"`
}

func main() {
	ctx := context.TODO()

	// Parse command line options.
	if _, err := flags.Parse(&opts); err != nil {
		log.Printf("[ERROR] Failed to parse flags: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Debug: %v\n", opts.Debug)

	tgBot, err := tgbotapi.NewBotAPI(opts.Telegram.Token)
	if err != nil {
		log.Fatalf("[ERROR] Can't make telegram bot, %v", err)
	}
	tgBot.Debug = opts.Debug

	log.Printf("Authorized on account %s", tgBot.Self.UserName)

	tgListener := events.TelegramListener{
		BotAPI: tgBot,
	}

	if err := tgListener.Do(ctx); err != nil {
		log.Fatalf("[ERROR] telegram listener failed, %v", err)
	}
}
