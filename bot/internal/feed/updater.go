package feedUpdater

import (
	"context"
	"fmt"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/rest"
	"log"
	"time"
)

const feedTTL = 2 * time.Minute

func StartFeedUpdater(client bot.Client) {
	ticker := time.NewTicker(feedTTL)
	go func() {
		for range ticker.C {
			updateFeeds(client)
		}
	}()
}

func updateFeeds(client bot.Client) {
	guilds, err := client.Rest().GetCurrentUserGuilds(context.TODO(), rest.GetCurrentUserGuildsOpts{})
	if err != nil {
		log.Fatalf("Fehler beim Abrufen der Guilds: %v", err)
	}

	// Ausgabe der Guild-IDs und -Namen
	for _, guild := range guilds {
		fmt.Printf("Guild ID: %s, Name: %s\n", guild.ID, guild.Name)
	}
}
