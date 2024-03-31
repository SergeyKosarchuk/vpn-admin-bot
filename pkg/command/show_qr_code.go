package command

import (
	"log"

	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ShowQRCode struct {
	Client client.APIClient
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

	data, err := c.Client.GetQRCode(id)

	if err != nil {
		return err
	}

	go func() {
		photo := tgbotapi.NewPhoto(output.ChatID, tgbotapi.FileBytes{Name: "qrcode.svg", Bytes: data})
		if _, err = c.Bot.Send(photo); err != nil {
			// TOTO: SVG files are not supported
			log.Println(err)
		}
	}()

	if err != nil {
		return err
	}

	output.Text = "QR Code will be send in a second."
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}
