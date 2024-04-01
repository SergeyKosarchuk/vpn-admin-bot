package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type List struct {
	Client APIClient
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

	if len(devices) > 0 {
		output.Text = listDevicesMsg(devices)
	} else {
		output.Text = "No devices."
	}

	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}
