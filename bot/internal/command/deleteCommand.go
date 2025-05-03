package command

import (
	"fmt"
	"github.com/disgoorg/disgo/events"
)

func HandleDeleteCommand(event *events.ApplicationCommandInteractionCreate) {
	fmt.Println("Delete Command Called")
}
