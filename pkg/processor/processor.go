package processor

import (
	"errors"

	botCommand "github.com/SergeyKosarchuk/vpn-admin-bot/pkg/command"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command interface {
	Prepare(output *tgbotapi.MessageConfig) error
	Action(input string, output *tgbotapi.MessageConfig) error
}

type MessageProcessor struct {
	selected       Command
	defaultCommand Command
	adminUsername  string
	commands       map[string]Command
}

func (mp *MessageProcessor) resposeToText(text string, output *tgbotapi.MessageConfig) error {
	if mp.selected != nil {
		selected := mp.selected
		// Command should be completed only once
		mp.selected = nil
		return selected.Action(text, output)
	}

	if mp.defaultCommand != nil {
		return mp.defaultCommand.Action(text, output)
	}

	return errors.New("default command is nil")
}

func (mp *MessageProcessor) resposeToCommand(commandName string, output *tgbotapi.MessageConfig) error {
	command, ok := mp.commands[commandName]

	if ok {
		// TODO: Do not select a command if user response is not required
		mp.selected = command
		return mp.selected.Prepare(output)
	}

	if mp.defaultCommand != nil {
		return mp.defaultCommand.Prepare(output)
	}

	return errors.New("default command is nil")
}

// Process message from user and prepare response
func (mp *MessageProcessor) MakeResponse(input tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	var err error
	response := tgbotapi.NewMessage(input.Chat.ID, "")

	if input.From.UserName != mp.adminUsername {
		response.Text = "You are not an admin."
		return response, err
	}

	if input.IsCommand() {
		err = mp.resposeToCommand(input.Command(), &response)
	} else {
		err = mp.resposeToText(input.Text, &response)
	}

	return response, err
}

// Create new MessageProcessor
func NewMessageProcessor(adminUsername string, client botCommand.APIClient, bot *tgbotapi.BotAPI) MessageProcessor {
	commands := make(map[string]Command)
	commands["ping"] = &botCommand.Ping{}
	commands["list"] = &botCommand.List{Client: client}
	commands["create"] = &botCommand.Create{Client: client}
	commands["enable"] = &botCommand.Enable{Client: client}
	commands["disable"] = &botCommand.Disable{Client: client}
	commands["delete"] = &botCommand.Delete{Client: client}
	commands["config"] = &botCommand.Config{Client: client, Bot: bot}
	commands["qrcode"] = &botCommand.ShowQRCode{Client: client, Bot: bot}
	processor := MessageProcessor{adminUsername: adminUsername, commands: commands, defaultCommand: &botCommand.EmptyCommand{}}

	return processor
}
