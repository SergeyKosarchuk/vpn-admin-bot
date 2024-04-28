package command

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skip2/go-qrcode"
)

type ShowQRCode struct {
	Client APIClient
	Bot    *tgbotapi.BotAPI
}

func (c *ShowQRCode) Prepare(output *tgbotapi.MessageConfig) error {
	devices, err := c.Client.List()
	if err != nil {
		return err
	}

	output.Text = "Select device to show qr code."
	output.ReplyMarkup = selectDeviceMarkup(devices)
	return nil
}

func (c *ShowQRCode) Action(input string, output *tgbotapi.MessageConfig) error {
	id, err := selectIdFromText(input)
	if err != nil {
		return err
	}

	data, err := c.Client.GetConfig(id)

	if err != nil {
		return err
	}

	go func() {
		// TODO: Wait for gorutine to complete
		png, err := qrcode.Encode(string(data), qrcode.Medium, 512)

		if err != nil {
			log.Println(err)
		} else {
			photo := tgbotapi.NewPhoto(output.ChatID, tgbotapi.FileBytes{Name: "qrcode.png", Bytes: png})

			if _, err = c.Bot.Send(photo); err != nil {
				log.Println(err)
			}
		}
	}()

	output.Text = "QR Code will be send in a second."
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}
