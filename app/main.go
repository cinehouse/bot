package main

import (
	"context"
	"fmt"
	"github.com/cinehouse/bot/app/bot"
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
	SysData string `long:"sys-data" env:"SYS_DATA" default:"data" description:"location of sys data"`
	Debug   bool   `long:"debug" description:"Show debug information"`
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

	multiBot := bot.MultiBot{}

	if sb, err := bot.NewSys(opts.SysData); err == nil {
		multiBot = append(multiBot, sb)
	} else {
		log.Printf("[ERROR] Failed to load sysbot, %v", err)
	}

	tgListener := events.TelegramListener{
		BotAPI: tgBot,
		Bots:   multiBot,
	}

	if err := tgListener.Do(ctx); err != nil {
		log.Fatalf("[ERROR] telegram listener failed, %v", err)
	}
}
