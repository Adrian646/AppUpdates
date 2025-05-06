package handler

import (
	"fmt"
	embedBuilder "github.com/Adrian646/AppUpdates/bot/internal/builder"
	apiclient "github.com/Adrian646/AppUpdates/bot/internal/service"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"os"
	"strings"
)

func HandleModal(e *events.ModalSubmitInteractionCreate) {
	parts := strings.Split(e.Data.CustomID, ":")

	switch parts[0] {
	case "register_app":
		platform := "android"
		if len(parts) > 1 {
			platform = parts[1]
		}
		handleRegisterApp(e, platform)
	case "delete_app":
		handleDeleteApp(e)
	default:
		fmt.Println("Unknown modal id: ", e.Data.CustomID)
	}
}

func handleRegisterApp(e *events.ModalSubmitInteractionCreate, platform string) {
	appID := e.Data.Text("app_id")
	client := apiclient.New(os.Getenv("API_BASE_URL"))

	feed, feedError := client.GetFeed(platform, appID)

	_, subscriptionError := client.CreateSubscription(e.GuildID().String(), e.Channel().ID().String(), platform, appID)

	if subscriptionError != nil {
		fmt.Println("Error creating subscription: ", subscriptionError)
		err := e.CreateMessage(discord.NewMessageCreateBuilder().
			SetEphemeral(true).
			AddEmbeds(embedBuilder.BuildErrorEmbed(
				"Error while creating subscription.\nPlease check the app ID and try again.\nIf this problem persists, please contact the developers.",
				nil,
				false,
			)).Build())
		if err != nil {
			return
		}
		return
	}

	messageError := e.CreateMessage(discord.NewMessageCreateBuilder().
		SetEphemeral(true).
		AddEmbeds(embedBuilder.BuildLoadingEmbed()).
		Build(),
	)
	if messageError != nil {
		return
	}

	if feedError != nil {
		fmt.Println("Error getting feed: ", feedError)
		err := e.CreateMessage(discord.NewMessageCreateBuilder().
			SetEphemeral(true).
			AddEmbeds(
				embedBuilder.BuildErrorEmbed(
					"Error while getting feed.\\nPlease check the app ID and try again.\\nIf this problem persists, please contact the developers.",
					nil,
					false,
				)).Build())
		if err != nil {
			return
		}
	}

	if platform == "android" {
		_, messageError = e.Client().Rest().CreateMessage(e.Channel().ID(), discord.NewMessageCreateBuilder().
			AddEmbeds(embedBuilder.BuildAndroidEmbed(feed)).
			Build(),
		)
		if messageError != nil {
			fmt.Println("Error sending channel message: ", messageError)
		}
	} else {
		_, messageError = e.Client().Rest().CreateMessage(e.Channel().ID(), discord.NewMessageCreateBuilder().
			AddEmbeds(embedBuilder.BuildIOSEmbed(feed)).
			Build(),
		)
		if messageError != nil {
			fmt.Println("Error sending channel message: ", messageError)
		}
	}

}

func handleDeleteApp(e *events.ModalSubmitInteractionCreate) {

}
