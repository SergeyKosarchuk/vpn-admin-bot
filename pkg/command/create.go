package command

import (
	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CreateCommand struct {
	Client client.APIClient
}

func (c *CreateCommand) Prepare(output *tgbotapi.MessageConfig) error {
	output.Text = "Input name for the new device."
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}

func (c *CreateCommand) Action(text string, output *tgbotapi.MessageConfig) error {
	err := c.Client.Create(text)
	if err != nil {
		return err
	}

	output.Text = "SUCCESS"
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}
