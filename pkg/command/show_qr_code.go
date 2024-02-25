package command

import (
	"io"
	"log"

	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


type ShowQRCode struct {
	Client client.APIClient
	Bot *tgbotapi.BotAPI
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

	dataReader, err := c.Client.GetQRCode(id)

	go func ()  {
		buffer, err := io.ReadAll(dataReader)

		if err != nil {
			log.Println("Unable to read a response into buffer", err)
			return
		}

		photo := tgbotapi.NewPhoto(output.ChatID, tgbotapi.FileBytes{Name: "qrcode.svg", Bytes: buffer})
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
