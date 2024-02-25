package processor

import (
	client "github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"
	botCommand "github.com/SergeyKosarchuk/vpn-admin-bot/pkg/command"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageProcessor struct {
	command botCommand.Command
	adminUsername string
	builder botCommand.CommandBuilder
}


func (mp *MessageProcessor) resposeToText(text string, output *tgbotapi.MessageConfig) error {
	return mp.command.Action(text, output)
}


func (mp *MessageProcessor) resposeToCommand(commandName string, output *tgbotapi.MessageConfig) error {
	mp.command = mp.builder.Build(commandName)
	return mp.command.Prepare(output)
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

func NewMessageProcessor(adminUsername string, client client.APIClient, bot *tgbotapi.BotAPI) *MessageProcessor {
	builder := botCommand.NewCommandBuilder(client, bot)

	return &MessageProcessor{
		command: builder.Build("empty"),
		adminUsername: adminUsername,
		builder: builder,
	}
}
