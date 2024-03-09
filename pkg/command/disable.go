package command

import (
	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Disable struct {
	Client client.APIClient
}

func (c *Disable) Prepare(output *tgbotapi.MessageConfig) error {
	devices, err := c.Client.List()
	if err != nil {
		return err
	}

	output.Text = "Select device to disable."
	output.ReplyMarkup = selectDeviceMarkup(devices)
	return nil
}

func (c *Disable) Action(input string, output *tgbotapi.MessageConfig) error {
	id, err := selectIdFromText(input)
	if err != nil {
		return err
	}

	err = c.Client.Disable(id)
	if err != nil {
		return err
	}

	output.Text = "ok"
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}
