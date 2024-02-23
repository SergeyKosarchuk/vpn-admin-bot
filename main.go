package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"
	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/manager"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)


func SelectDeviceMarkup(devices []client.Device) tgbotapi.ReplyKeyboardMarkup {
	rows := make([][]tgbotapi.KeyboardButton, len(devices))

	for idx, device := range devices {
		text := fmt.Sprintf("%d. %s", idx, device.Name)
		row := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(text))
		rows[idx] = row
	}

	return tgbotapi.NewReplyKeyboard(rows...)
}


func SelectDeviceMarkupInline(devices []client.Device) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, len(devices))

	for idx, device := range devices {
		text := fmt.Sprintf("%d. %s", idx, device.Name)
		button := tgbotapi.NewInlineKeyboardButtonData(text, device.Id)
		rows[idx] = tgbotapi.NewInlineKeyboardRow(button)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}


func CreateDeviceMessage(device client.Device, chatID int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, device.Name)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("More", "more." + device.Id),
	),
    tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("✅", "enable." + device.Id),
        tgbotapi.NewInlineKeyboardButtonData("🗑", "delete." + device.Id),
				tgbotapi.NewInlineKeyboardButtonData("⚙", "config." + device.Id),
    ),
	)
	return msg
}



func ListDevicesMsg(devices []client.Device) string {
	var sb strings.Builder

	for idx, device := range devices {
		sb.WriteString(fmt.Sprintf("%d. %s\n", idx, device))
	}

	return sb.String()
}



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
	manager := manager.DeviceManager{Client: wgClient}
	err = manager.Fetch()
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
		if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
					panic(err)
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
					panic(err)
			}
		}

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
			msg.Text = "PONG"
		case "enable":
			msg.Text = "Select device number to enable"
			msg.ReplyMarkup = SelectDeviceMarkupInline(manager.Devices)
		case "disable":
			msg.Text = "Select device number to enable"
			msg.ReplyMarkup = SelectDeviceMarkupInline(manager.Devices)
		case "list":
				msg.Text = ListDevicesMsg(manager.Devices)
		case "close":
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "new":
				for _, device := range manager.Devices {
					msg := CreateDeviceMessage(device, update.Message.Chat.ID)
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
				}
				continue
		default:
			msg.Text = "I don't know that command"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
