package main

import (
	"os"

	_ "./trivia"
	_ "./trumpisms"
	"github.com/go-chat-bot/bot/slack"
	_ "github.com/go-chat-bot/plugins/encoding"
	// _ "github.com/go-chat-bot/plugins/catgif"
	// _ "github.com/go-chat-bot/plugins/chucknorris"
	// Import all the commands you wish to use
)

func main() {
	slack.Run(os.Getenv("kacibot"))
}
