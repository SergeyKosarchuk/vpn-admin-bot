package command

import (
	"io"
	"log"

	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	Client client.APIClient
	Bot    *tgbotapi.BotAPI
}

func (c *Config) Prepare(output *tgbotapi.MessageConfig) error {
	devices, err := c.Client.List()
	if err != nil {
		return err
	}

	output.Text = "Select device to get config."
	output.ReplyMarkup = selectDeviceMarkup(devices)
	return nil
}

func (c *Config) Action(input string, output *tgbotapi.MessageConfig) error {
	id, err := selectIdFromText(input)
	if err != nil {
		return err
	}

	dataReader, err := c.Client.GetConfig(id)

	go func() {
		buffer, err := io.ReadAll(dataReader)
		if err != nil {
			log.Println("Unable to read a response into buffer", err)
			return
		}

		photo := tgbotapi.NewDocument(output.ChatID, tgbotapi.FileBytes{Name: "vpn.conf", Bytes: buffer})
		if _, err = c.Bot.Send(photo); err != nil {
			log.Println(err)
		}
	}()

	if err != nil {
		return err
	}

	output.Text = "Ok"
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}
