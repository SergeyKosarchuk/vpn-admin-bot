package command

import (
	"fmt"
	"strings"

	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func selectDeviceMarkup(devices []client.DeviceResponse) tgbotapi.ReplyKeyboardMarkup {
	rows := make([][]tgbotapi.KeyboardButton, len(devices))

	for idx, device := range devices {
		text := fmt.Sprintf("%d. %s (%s)\n", idx, device.Id, device.Name)
		row := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(text))
		rows[idx] = row
	}

	return tgbotapi.NewReplyKeyboard(rows...)
}

func listDevicesMsg(devices []client.DeviceResponse) string {
	var sb strings.Builder

	for idx, device := range devices {
		sb.WriteString(fmt.Sprintf("%d. %s (%s)\n", idx, device.Id, device.Name))
	}

	return sb.String()
}


func selectIdFromText(text string) (string, error) {
	parts := strings.Split(text, " ")

	if len(parts) != 3 {
		return "", fmt.Errorf("unable to select an id from `%s` string", text)
	}

	return parts[1], nil
}
