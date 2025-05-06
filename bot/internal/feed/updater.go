package feedUpdater

import (
	"fmt"
	embedBuilder "github.com/Adrian646/AppUpdates/bot/internal/builder"
	api "github.com/Adrian646/AppUpdates/bot/internal/service"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"os"
	"time"
)

const feedTTL = 3 * time.Minute

func StartFeedUpdater(client bot.Client) {
	ticker := time.NewTicker(feedTTL)
	go func() {
		for range ticker.C {
			updateFeeds(client)
		}
	}()
}

func updateFeeds(client bot.Client) {
	service := api.New(os.Getenv("API_BASE_URL"))

	updates, err := service.GetFeedUpdates()
	if err != nil {
		return
	}

	for _, sub := range updates {
		feed := &sub.AppFeed
		channelID := snowflake.MustParse(sub.ChannelID)
		fmt.Printf("Updating feed %s for channel %s\n", feed.AppID, channelID.String())
		if feed.Platform == "android" {
			_, messageError := client.Rest().CreateMessage(channelID, discord.NewMessageCreateBuilder().
				AddEmbeds(embedBuilder.BuildAndroidEmbed(feed)).
				Build(),
			)
			if messageError != nil {
				fmt.Println("Error sending channel message: ", messageError)
			}
		} else {
			_, messageError := client.Rest().CreateMessage(channelID, discord.NewMessageCreateBuilder().
				AddEmbeds(embedBuilder.BuildIOSEmbed(feed)).
				Build(),
			)
			if messageError != nil {
				fmt.Println("Error sending channel message: ", messageError)
			}
		}
	}
}
