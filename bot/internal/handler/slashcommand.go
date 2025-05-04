package handler

import (
	"github.com/Adrian646/AppUpdates/bot/internal/command"
	"github.com/disgoorg/disgo/events"
)

func HandleSlashCommand(e *events.ApplicationCommandInteractionCreate) {
	switch e.Data.CommandName() {
	case "register":
		command.HandleRegisterCommand(e)
	case "delete":
		command.HandleDeleteCommand(e)
	}
}
