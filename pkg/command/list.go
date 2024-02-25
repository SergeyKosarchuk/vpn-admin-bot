package command

import (
	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type List struct {
	Client client.APIClient
}


func (c *List) Action(input string, output *tgbotapi.MessageConfig) error {
	output.Text = "Show menu"
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}


func (c *List) Prepare(output *tgbotapi.MessageConfig) error {
	devices, err := c.Client.List()

		if err != nil {
			return err
		}

		output.Text = listDevicesMsg(devices)
		output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}
