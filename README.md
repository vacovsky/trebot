# Trebot - trivia plugin

#### for go-chat-bot https://github.com/go-chat-bot/bot

## Usage / Bot Interaction

### Get a new question

```slack
!trivia new
```

<img src="https://raw.githubusercontent.com/vacoj/trebot/master/img/s2.png">

### Answer a question

```slack
!trivia answer pineapple
```

<img src="https://raw.githubusercontent.com/vacoj/trebot/master/img/s1.png">

### View the Scoreboard

```slack
!trivia scoreboard
```

```
+------+------------------+--------+---------+-------+-----+----------+
| RANK |       USER       | SCORE  | CORRECT | WRONG | NEW | ACCURACY |
+------+------------------+--------+---------+-------+-----+----------+
|    1 | casey            | 225300 |     164 |   122 |  15 |    0.573 |
|    2 | viki             |  63800 |     111 |   234 |  42 |    0.322 |
|    3 | vacoj            |  60700 |      40 |    65 |  16 |    0.381 |
|    4 | bologna          |  58000 |      61 |   128 |  50 |    0.323 |
|    5 | josh             |  55800 |      53 |   310 |  45 |    0.146 |
|    6 | krem             |  38200 |      73 |   196 |  33 |    0.271 |
|    7 | joey             |  18000 |       2 |     2 |   0 |    0.500 |
|    8 | k-w              |   1000 |       0 |     0 |   0 |    0.000 |
+------+------------------+--------+---------+-------+-----+----------+
```

### View "About" information

```slack
!trivia about
```

## Setup / Installation

### Stand Alone

``` bash
go get -u "github.com/vacoj/trebot"
go get -u "github.com/vacoj/trebot/trivia"
go get -u "github.com/go-chat-bot/bot/slack"
# Set your bot token for slack as an environment variable called "trebot"
export trebot="xxxx-yourslackbotkey"
trebot
```

### As a plugin for an existing bot

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