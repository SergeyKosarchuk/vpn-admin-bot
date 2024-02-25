package command

import (
	api "github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


type Command interface {
	Prepare(output *tgbotapi.MessageConfig) error
	Action(input string, output *tgbotapi.MessageConfig) error
}


type CommandBuilder interface {
	Build(name string) Command
}

type commandBuilder struct {
	Client api.APIClient
	Bot *tgbotapi.BotAPI
}

func (cb commandBuilder) Build(name string) Command {
	switch name {
	case "ping":
		return &PingCommand{}
	case "list":
		return &List{Client: cb.Client}
	case "create":
		return &CreateCommand{Client: cb.Client}
	case "enable":
		return &Enable{Client: cb.Client}
	case "disable":
		return &Disable{Client: cb.Client}
	case "delete":
		return &Delete{Client: cb.Client}
	case "code":
		return &ShowQRCode{Client: cb.Client, Bot: cb.Bot}
	case "config":
		return &Config{Client: cb.Client, Bot: cb.Bot}
	default:
		return &EmptyCommand{}
	}
}

func NewCommandBuilder(client api.APIClient, bot *tgbotapi.BotAPI) CommandBuilder {
	return &commandBuilder{Client: client, Bot: bot}
}
