package main

import (
	"os"

	"github.com/go-chat-bot/bot/slack"
	_ "github.com/vacoj/trebot/trivia"
)

func main() {
	slack.Run(os.Getenv("kacibot"))
}
