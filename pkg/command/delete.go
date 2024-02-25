package command

import (
	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


type Delete struct {
	Client client.APIClient
}


func (c *Delete) Prepare(output *tgbotapi.MessageConfig) error {
	devices, err := c.Client.List()

	if err != nil {
		return err
	}

	output.Text = "Select device to delete."
	output.ReplyMarkup = selectDeviceMarkup(devices)
	return nil
}

func (c *Delete) Action(input string, output *tgbotapi.MessageConfig) error {
	id, err := selectIdFromText(input)

	if err != nil {
		return err
	}

	err = c.Client.Delete(id)

	if err != nil {
		return err
	}

	output.Text = "ok"
	output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return nil
}
