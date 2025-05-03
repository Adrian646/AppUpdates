package main

import (
	"context"
	"fmt"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/json"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:                     "register",
		Description:              "Register a new app that should be tracked",
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionAdministrator),
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "platform",
				Description: "Choose the platform",
				Required:    true,
				Choices: []discord.ApplicationCommandOptionChoiceString{
					{Name: "Android", Value: "android"},
					{Name: "IOS", Value: "ios"},
				},
			},
		},
	},
	discord.SlashCommandCreate{
		Name:                     "delete",
		Description:              "Delete an app that is currently being tracked",
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionAdministrator),
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "platform",
				Description: "Choose the platform",
				Required:    true,
				Choices: []discord.ApplicationCommandOptionChoiceString{
					{Name: "Android", Value: "android"},
					{Name: "IOS", Value: "ios"},
				},
			},
		},
	},
}

func main() {
	fmt.Println("Starting bot...")

	err := godotenv.Load("../.env")

	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	client, err := disgo.New(os.Getenv("BOT_TOKEN"),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
			),
		),
	)

	if err != nil {
		panic(err)
	}

	if _, err = client.Rest().SetGlobalCommands(client.ApplicationID(), commands); err != nil {
		panic(err)
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		panic(err)
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}
