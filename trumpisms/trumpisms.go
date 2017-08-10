package trumpisms

import (
	"fmt"

	"github.com/go-chat-bot/bot"
)

func trumpify(command *bot.Cmd) (string, error) {

	result, _ := trumpIt()
	return result, nil
}

func trumpIt() (string, error) {
	return "test", nil
}

func init() {
	fmt.Println("Registering Trumpisms...")
	bot.RegisterCommand(
		"trumpify",
		"Turn your phrase into a Trumpism!",
		"{phrase}",
		trumpify)
}
