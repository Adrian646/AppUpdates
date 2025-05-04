package command

import (
	"fmt"
	"github.com/disgoorg/disgo/events"
)

func HandleDeleteCommand(e *events.ApplicationCommandInteractionCreate) {
	fmt.Println("Delete Command Called")
}
