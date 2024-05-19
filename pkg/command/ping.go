package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Ping struct{}

func (c *Ping) Action(input string, output *tgbotapi.MessageConfig) error {
	output.Text = "Show menu"
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}

func (c *Ping) Prepare(output *tgbotapi.MessageConfig) error {
	output.Text = "PONG"
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}
