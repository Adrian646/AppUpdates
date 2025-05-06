package command

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func HandleRegisterCommand(e *events.ApplicationCommandInteractionCreate) {
	platform := e.SlashCommandInteractionData().String("platform")

	modal := discord.NewModalCreateBuilder().
		SetCustomID(fmt.Sprintf("register_app:%s", platform)).
		SetTitle(fmt.Sprintf("Register an %s app", platform)).
		AddActionRow(discord.NewTextInput("app_id", discord.TextInputStyleShort, "Please enter your app id")).
		Build()

	err := e.Modal(modal)
	if err != nil {
		return
	}
}
