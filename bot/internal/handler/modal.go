package handler

import (
	"fmt"
	apiclient "github.com/Adrian646/AppUpdates/bot/internal/service"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"os"
)

func HandleModal(e *events.ModalSubmitInteractionCreate) {
	switch e.Data.CustomID {
	case "register_app":
		handleRegisterApp(e)
	case "delete_app":
		handleDeleteApp(e)
	default:
		fmt.Println("Unknown modal id: ", e.Data.CustomID)
	}
}

func handleRegisterApp(e *events.ModalSubmitInteractionCreate) {
	appID := e.Data.Text("app_id")
	client := apiclient.New(os.Getenv("API_BASE_URL"))
	embed := discord.NewEmbedBuilder()

	feed, err := client.GetFeed("android", appID)
	if err != nil {
		fmt.Println("Error getting feed: ", err)
		embed.SetColor(0xff0000)
		embed.SetDescription("Error while getting feed.\nPlease check the app ID and try again.\nIf this problem persists, please contact the developers.")
		err := e.CreateMessage(discord.NewMessageCreateBuilder().AddEmbeds(embed.Build()).Build())
		if err != nil {
			return
		}
	}

	_, err = client.CreateSubscription(e.GuildID().String(), e.Channel().ID().String(), "android", appID)

	if err != nil {
		fmt.Println("Error creating subscription: ", err)
		embed.SetColor(0xff0000)
		embed.SetDescription("Error while creating subscription.\nPlease check the app ID and try again.\nIf this problem persists, please contact the developers.")
		err := e.CreateMessage(discord.NewMessageCreateBuilder().AddEmbeds(embed.Build()).Build())
		if err != nil {
			return
		}
		return
	}

	embed.SetColor(0x00ff00)

	embed.SetAuthor("Android Update", "", "https://upload.wikimedia.org/wikipedia/commons/d/d7/Android_robot.svg")

	embed.SetTitlef("%s v%s is available!", feed.AppID, feed.Version)
	embed.SetURLf("https://play.google.com/store/apps/details?id=%s", feed.AppID)
	embed.SetThumbnail(feed.AppIconURL)

	embed.AddField("**Dynamic Details**",
		fmt.Sprintf(
			"**App ID:** `%s`\n**Developer:** %s\n**Updated on:** %s\n**Downloads:** %s+",
			feed.AppID,
			feed.Developer,
			feed.UpdatedOn.Format("2006-01-02"),
			feed.DownloadCount,
		),
		false)

	embed.AddField("**Release Notes**", feed.ReleaseNotes, false)

	if feed.AppBannerURL != "" {
		embed.SetImage(feed.AppBannerURL)
	}

	err = e.CreateMessage(discord.NewMessageCreateBuilder().
		AddEmbeds(embed.Build()).
		Build(),
	)
	if err != nil {
		return
	}

}

func handleDeleteApp(e *events.ModalSubmitInteractionCreate) {

}
