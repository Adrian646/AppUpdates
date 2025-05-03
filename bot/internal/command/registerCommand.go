package command

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func HandleRegisterCommand(event *events.ApplicationCommandInteractionCreate) {
	platform := event.SlashCommandInteractionData().String("platform")

	modal := discord.NewModalCreateBuilder().
		SetCustomID("register_app").
		SetTitle(fmt.Sprintf("Register an %s app", platform)).
		AddActionRow(discord.NewTextInput("app_id", discord.TextInputStyleShort, "Please enter your app id")).
		Build()

	err := event.Modal(modal)
	if err != nil {
		return
	}
}
