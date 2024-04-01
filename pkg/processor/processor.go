package processor

import (
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
		return mp.selected.Action(text, output)
	}

	return mp.defaultCommand.Action(text, output)
}

func (mp *MessageProcessor) resposeToCommand(commandName string, output *tgbotapi.MessageConfig) error {
	command, ok := mp.commands[commandName]

	if ok {
		mp.selected = command
		return mp.selected.Prepare(output)
	}

	return mp.defaultCommand.Prepare(output)
	
}

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

func NewMessageProcessor(adminUsername string, client botCommand.APIClient, bot *tgbotapi.BotAPI) MessageProcessor {
	commands := make(map[string]Command)
	commands["ping"] = &botCommand.PingCommand{}
	commands["list"] = &botCommand.List{Client: client}
	commands["create"] = &botCommand.CreateCommand{Client: client}
	commands["enable"] = &botCommand.Enable{Client: client}
	commands["disable"] = &botCommand.Disable{Client: client}
	commands["delete"] = &botCommand.Delete{Client: client}
	commands["config"] = &botCommand.Config{Client: client, Bot: bot}
	processor := MessageProcessor{adminUsername: adminUsername, commands: commands, defaultCommand: &botCommand.EmptyCommand{}}

	return processor
}
