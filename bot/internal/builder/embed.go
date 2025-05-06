package embedBuilder

import (
	"fmt"
	api "github.com/Adrian646/AppUpdates/bot/internal/service"
	"github.com/disgoorg/disgo/discord"
)

func BuildAndroidEmbed(feed *api.AppFeed) discord.Embed {
	embed := discord.NewEmbedBuilder()

	embed.SetColor(0x00ff00)

	embed.SetAuthor("Android Update", "", "https://cdn.discordapp.com/emojis/1369094424474423439.webp")

	embed.SetTitlef("%s v%s is available!", feed.AppName, feed.Version)
	embed.SetURLf("https://play.google.com/store/apps/details?id=%s", feed.AppID)
	embed.SetThumbnail(feed.AppIconURL)

	embed.AddField("**Dynamic Details**",
		fmt.Sprintf(
			"**App ID:** `%s`\n**Developer:** %s\n**Updated on:** %s\n**Downloads:** %s",
			feed.AppID,
			feed.Developer,
			feed.UpdatedOn.Format("2006-01-02"),
			feed.DownloadCount,
		),
		false)

	if feed.ReleaseNotes != "" {
		embed.AddField("**Release Notes**", feed.ReleaseNotes, false)
	}

	if feed.AppBannerURL != "" {
		embed.SetImage(feed.AppBannerURL)
	}

	return embed.Build()
}

func BuildIOSEmbed(feed *api.AppFeed) discord.Embed {
	embed := discord.NewEmbedBuilder()

	embed.SetColor(0x1aa9f9)

	embed.SetAuthor("IOS Update", "", "https://cdn.discordapp.com/emojis/1369094367800852570.webp")

	embed.SetTitlef("%s v%s is available!", feed.AppName, feed.Version)
	embed.SetURLf("https://play.google.com/store/apps/details?id=%s", feed.AppID)
	embed.SetThumbnail(feed.AppIconURL)

	embed.AddField("**Dynamic Details**",
		fmt.Sprintf(
			"**App ID:** `%s`\n**Developer:** %s\n**Updated on:** %s",
			feed.AppID,
			feed.Developer,
			feed.UpdatedOn.Format("2006-01-02"),
		),
		false)

	if feed.ReleaseNotes != "" {
		embed.AddField("**Release Notes**", feed.ReleaseNotes, false)
	}

	if feed.AppBannerURL != "" {
		embed.SetImage(feed.AppBannerURL)
	}

	return embed.Build()
}

func BuildLoadingEmbed() discord.Embed {
	embed := discord.NewEmbedBuilder()

	embed.SetColor(0xFFBF00)
	embed.SetDescriptionf("%s Fetching app data...\nThis may take a few seconds. Please be patient!", "<a:loading:1369093708594811003>")

	return embed.Build()
}

func BuildErrorEmbed(msg string, err error, showError bool) discord.Embed {
	embed := discord.NewEmbedBuilder()

	embed.SetColor(0xff0000)

	if showError {
		embed.SetDescriptionf("%s: %s", msg, err.Error())
	} else {
		embed.SetDescriptionf("%s", msg)
	}

	return embed.Build()
}
