package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


type PingCommand struct {}

func (c *PingCommand) Action(input string, output *tgbotapi.MessageConfig) error {
	output.Text = "Show menu"
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}

func (c *PingCommand) Prepare(output *tgbotapi.MessageConfig) error {
	output.Text = "PONG"
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}
