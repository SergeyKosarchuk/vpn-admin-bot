package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	adminUsername := os.Getenv("VPN_ADMIN_TELEGRAM_USERNAME")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	log.Printf("Receiving commands from %s", adminUsername)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message == nil { // ignore any non-Message updates
				continue
			}
			if !update.Message.IsCommand() { // ignore any non-command Messages
				continue
			}
			if update.Message.From.UserName != adminUsername { // accept commands only from admin
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "ping":
					msg.Text = "pong"
			default:
					msg.Text = "I don't know that command"
			}

			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}
