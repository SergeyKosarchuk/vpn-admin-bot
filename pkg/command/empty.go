package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type EmptyCommand struct{}

func (c *EmptyCommand) Prepare(output *tgbotapi.MessageConfig) error {
	output.Text = "Please select a command"
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}

func (c *EmptyCommand) Action(text string, output *tgbotapi.MessageConfig) error {
	output.Text = "Unknown"
	return nil
}
