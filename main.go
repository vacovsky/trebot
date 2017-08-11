package main

import (
	"os"

	_ "./trivia"
	"github.com/go-chat-bot/bot/slack"
)

func main() {
	slack.Run(os.Getenv("kacibot"))
}
