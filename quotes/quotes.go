package quotes

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-chat-bot/bot"
)

func quote(command *bot.Cmd) (string, error) {

	result, _ := getQuote()
	return result, nil
}

func getQuote() (string, error) {
	ind := random(0, len(quotes))
	return string(quotes[ind]), nil
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func init() {
	fmt.Println("Registering InfraQuotes...")
	bot.RegisterCommand(
		"quote",
		"Displays a priceless quote",
		"No params.",
		quote)
}
