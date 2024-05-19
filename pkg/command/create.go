package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Create struct {
	Client APIClient
}

func (c *Create) Prepare(output *tgbotapi.MessageConfig) error {
	output.Text = "Input name for the new device."
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}

func (c *Create) Action(text string, output *tgbotapi.MessageConfig) error {
	err := c.Client.Create(text)
	if err != nil {
		return err
	}

	output.Text = "SUCCESS"
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}
