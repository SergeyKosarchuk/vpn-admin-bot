package main

import (
	"log"
	"os"

	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"

	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/processor"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)


func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
	telegramToken, ok := os.LookupEnv("TELEGRAM_API_TOKEN")
	if !ok {
		log.Fatalf("telegram token not set")
	}

	adminUsername, ok := os.LookupEnv("TELEGRAM_USERNAME")
	if !ok {
		log.Fatalf("adminUsername token not set")
	}

	adminUrl, ok := os.LookupEnv("VPN_ADMIN_URL")
	if !ok {
		log.Fatalf("adminUrl token not set")
	}

	adminPassword, ok := os.LookupEnv("VPN_ADMIN_PASSWORD")
	if !ok {
		log.Fatalf("adminUrl token not set")
	}

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}
	wgClient, err := client.NewWGClient(adminUrl, adminPassword)
	if err != nil {
		log.Panic(err)
	}
	mp := processor.NewMessageProcessor(adminUsername, wgClient)

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
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		msg, err := mp.MakeResponse(*update.Message)

		if err != nil {
			log.Printf("Error %s", err)

			if msg.Text == "" {
				msg.Text = "Unknown Error"
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			}
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
