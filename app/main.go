package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-pkgz/lgr"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jessevdk/go-flags"

	"github.com/cinehouse/bot/app/bot"
	"github.com/cinehouse/bot/app/events"
	"github.com/cinehouse/bot/app/reporter"
)

var opts struct {
	Telegram struct {
		Token string `long:"token" env:"TOKEN" description:"Telegram bot token"`
	} `group:"telegram" namespace:"telegram" env-namespace:"TELEGRAM"`
	LogsPath   string           `short:"l" long:"logs" env:"TELEGRAM_LOGS" default:"logs" description:"path to logs"`
	SuperUsers events.SuperUser `long:"super" description:"super-users"`
	SysData    string           `long:"sys-data" env:"SYS_DATA" default:"data" description:"location of sys data"`
	Debug      bool             `long:"debug" description:"Show debug information"`
}

func main() {
	ctx := context.TODO()

	// Parse command line options.
	if _, err := flags.Parse(&opts); err != nil {
		log.Printf("[ERROR] Failed to parse flags: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Debug: %v\n", opts.Debug)
	setupLog(opts.Debug)
	log.Printf("[INFO] super users: %v", opts.SuperUsers)

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
		TgBotAPI:      tgBot,
		MessageLogger: reporter.NewLogger(opts.LogsPath),
		Bots:          multiBot,
	}

	if err := tgListener.Do(ctx); err != nil {
		log.Fatalf("[ERROR] telegram listener failed, %v", err)
	}
}

func setupLog(debug bool) {
	logOpts := []lgr.Option{lgr.Msec, lgr.LevelBraces}
	if debug {
		logOpts = []lgr.Option{lgr.Debug, lgr.CallerFile, lgr.CallerFunc, lgr.Msec, lgr.LevelBraces}
	}
	lgr.SetupStdLogger(logOpts...)
}
