# Trebot - trivia plugin

#### for go-chat-bot https://github.com/go-chat-bot/bot

![Circle CI](https://circleci.com/gh/vacoj/trebot/tree/master.svg?style=svg)

## Usage / Bot Interaction

### Get a new question

```slack
!trivia new
```

<img src="https://raw.githubusercontent.com/vacoj/trebot/master/img/s2.png">

### Answer a question

```slack
!trivia answer pineapple
or
pineapple
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
export TREBOT_KEY="xxxx-yourslackbotkey"
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
	slack.Run(os.Getenv("TREBOT_KEY"))
}
```

## Contributing

Pull requests are encouraged. Please submit unit tests with any submitted patch. Also, please make it clear what your patch does in branch and commit messages.

1. Fork the repository
2. Create your feature branch `git checkout -b change-to-thing`
3. Commit your changes `git commit -am 'changes to the thing'`
4. Push to the branch `git push origin change-to-thing`
5. Create new pull request
6. Be patient

