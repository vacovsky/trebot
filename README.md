# Trebot - trivia plugin
## for go-chat-bot https://github.com/go-chat-bot/bot

### Setup - stand alone

``` bash
go get -u "github.com/vacoj/trebot"
# Set your bot token for slack as an environment variable called "trebot"
export trebot="xxxx-yourslackbotkey"
trebot
```

### Setup - plugin for existing bot

#### Install plugin

``` bash
go get -u "github.com/vacoj/trebot/trivia"
#Set your bot token for slack as an environment variable called "trebot"
export trebot="xxxx-yourslackbotkey"
```

#### Update your bot main to include the plugin

``` go
package main

import (
	"os"

	_ "github.com/vacoj/trebot/trivia"
	"github.com/go-chat-bot/bot/slack"
)

func main() {
	slack.Run(os.Getenv("trebot"))
}
```