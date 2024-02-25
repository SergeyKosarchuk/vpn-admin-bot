package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


type PingCommand struct {NoInputCommand}


func (c *PingCommand) Prepare(output *tgbotapi.MessageConfig) error {
	output.Text = "PONG"
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}
