package processor

import (
	"fmt"
	"strings"

	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageProcessor struct {
	client client.APIClient
	command string
	adminUsername string
}

func SelectDeviceMarkup(devices []client.DeviceResponse) tgbotapi.ReplyKeyboardMarkup {
	rows := make([][]tgbotapi.KeyboardButton, len(devices))

	for idx, device := range devices {
		text := fmt.Sprintf("%d. %s (%s)\n", idx, device.Id, device.Name)
		row := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(text))
		rows[idx] = row
	}

	return tgbotapi.NewReplyKeyboard(rows...)
}


func ListDevicesMsg(devices []client.DeviceResponse) string {
	var sb strings.Builder

	for idx, device := range devices {
		sb.WriteString(fmt.Sprintf("%d. %s (%s)\n", idx, device.Id, device.Name))
	}

	return sb.String()
}


const NOT_AN_ADMIN = "You are not an admin."
const SELECT_TO_DELETE = "Select device id to delete."
const SELECT_TO_ENABLE = "Select device id to enable."
const SELECT_TO_DISABLE = "Select device id to disable."
const UNKNOWN_COMMAND = "Unknown command."
const PONG = "PONG"
const SUCCESS = "Success"


const CREATE_COMMAND = "CREATE"
const DELETE_COMMAND = "DELETE"
const ENABLE_COMMAND = "ENABLE"
const DISABLE_COMMAND = "DISABLE"
const EMPTY_COMAND = ""


func (mp *MessageProcessor) resposeToCommand(command string, output *tgbotapi.MessageConfig) error {
	switch command {
	case "ping":
		output.Text = PONG
	case "list":
		devices, err := mp.client.List()

		if err != nil {
			return err
		}

		output.Text = ListDevicesMsg(devices)
		output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		mp.command = EMPTY_COMAND
	case "delete":
		devices, err := mp.client.List()

		if err != nil {
			return err
		}

		output.Text = SELECT_TO_DELETE
		output.ReplyMarkup = SelectDeviceMarkup(devices)
		mp.command = DELETE_COMMAND
	case "create":
		output.Text = "Input name for the new device."
		mp.command = CREATE_COMMAND
	case "enable":
		devices, err := mp.client.List()

		if err != nil {
			return err
		}

		output.Text = SELECT_TO_ENABLE
		output.ReplyMarkup = SelectDeviceMarkup(devices)
		mp.command = ENABLE_COMMAND
	case "disable":
	devices, err := mp.client.List()

		if err != nil {
			return err
		}

		output.Text = SELECT_TO_DISABLE
		output.ReplyMarkup = SelectDeviceMarkup(devices)
		mp.command = DISABLE_COMMAND
	default:
		output.Text = UNKNOWN_COMMAND
		mp.command = EMPTY_COMAND
		output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}
	return nil
}

func (mp *MessageProcessor) resposeToText(text string, output *tgbotapi.MessageConfig) error {
	if mp.command == "" {
		output.Text = "No command"
		return nil
	}
	switch mp.command {
	case CREATE_COMMAND:
		err := mp.client.Create(text)

		if err != nil {
			return nil
		}

		output.Text = SUCCESS
		output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}
	output.Text = "NOT IMP"
	return nil
}

func (mp *MessageProcessor) MakeResponse(input tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	var err error
	response := tgbotapi.NewMessage(input.Chat.ID, "")

	if input.From.UserName != mp.adminUsername {
		response.Text = NOT_AN_ADMIN
		return response, err
	}

	if input.IsCommand() {
		err = mp.resposeToCommand(input.Command(), &response)
	} else {
		err = mp.resposeToText(input.Text, &response)
	}

	return response, err
}

func NewMessageProcessor(adminUsername string, client client.APIClient) *MessageProcessor {
	return &MessageProcessor{
		client: client,
		command: EMPTY_COMAND,
		adminUsername: adminUsername,
	}
}
