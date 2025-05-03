package handler

import (
	"github.com/Adrian646/AppUpdates/bot/internal/command"
	"github.com/disgoorg/disgo/events"
)

func HandleSlashCommand(event *events.ApplicationCommandInteractionCreate) {
	switch event.Data.CommandName() {
	case "register":
		command.HandleRegisterCommand(event)
	case "delete":
		command.HandleDeleteCommand(event)
	}
}
