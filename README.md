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

## Blinky Lights

```

1516838759.729986 [0 129.65.184.221:55828] "AUTH" "fbn4u8ow0bf7w389w4780pfobawuy8b4w780opfbaewu8faopghu9wopfbesiu"
1516838759.731502 [0 129.65.184.221:55828] "PUBLISH" "ws2812b_0" "{'senselight': 0, 'brightness': 10, 'color': ['tan'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 630}"
1516838760.040196 [0 129.65.184.221:55828] "PUBLISH" "ws2812b_0" "{'senselight': 0, 'brightness': 10, 'color': ['palegoldenrod'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 630}"
1516838760.343341 [0 129.65.184.221:55828] "PUBLISH" "ws2812b_0" "{'senselight': 0, 'brightness': 10, 'color': ['firebrick'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 630}"
1516838760.646675 [0 129.65.184.221:55828] "PUBLISH" "ws2812b_0" "{'senselight': 0, 'brightness': 10, 'color': ['plum'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 630}"
1516838760.950227 [0 129.65.184.221:55828] "PUBLISH" "ws2812b_0" "{'senselight': 0, 'brightness': 10, 'color': ['darkslategray'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 630}"
1516838835.413406 [0 129.65.184.221:55830] "AUTH" "fbn4u8ow0bf7w389w4780pfobawuy8b4w780opfbaewu8faopghu9wopfbesiu"
1516838835.414897 [0 129.65.184.221:55830] "PUBLISH" "ws2812b_0" "{'senselight': 0, 'brightness': 10, 'color': ['lightgoldenrodyellow'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 630}"
1516838835.723626 [0 129.65.184.221:55830] "PUBLISH" "ws2812b_0" "{'senselight': 0, 'brightness': 10, 'color': ['seashell'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 630}"
1516838836.027203 [0 129.65.184.221:55830] "PUBLISH" "ws2812b_0" "{'senselight': 0, 'brightness': 10, 'color': ['seashell'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 630}"
1516838836.330787 [0 129.65.184.221:55830] "PUBLISH" "ws2812b_0" "{'senselight': 0, 'brightness': 10, 'color': ['darksalmon'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 630}"
1516838836.634226 [0 129.65.184.221:55830] "PUBLISH" "ws2812b_0" "{'senselight': 0, 'brightness': 10, 'color': ['peachpuff'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 630}"
1516838838.803780 [0 129.65.184.221:55832] "AUTH" "fbn4u8ow0bf7w389w4780pfobawuy8b4w780opfbaewu8faopghu9wopfbesiu"
1516838838.805286 [0 129.65.184.221:55832] "PUBLISH" "ws2812b_0" "{'senselight': 0, 'brightness': 10, 'color': ['palevioletred'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 630}"
1516838839.113754 [0 129.65.184.221:55832] "PUBLISH" "ws2812b_0" "{'senselight': 0, 'brightness': 10, 'color': ['paleturquoise'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 630}"
1516838839.417192 [0 129.65.184.221:55832] "PUBLISH" "ws2812b_0" "{'senselight': 0, 'brightness': 10, 'color': ['skyblue'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 630}"
1516838839.720511 [0 129.65.184.221:55832] "PUBLISH" "ws2812b_0" "{'senselight': 0, 'brightness': 10, 'color': ['lightsteelblue'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 630}"
1516838840.023737 [0 129.65.184.221:55832] "PUBLISH" "ws2812b_0" "{'senselight': 0, 'brightness': 10, 'color': ['floralwhite'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 630}"

```

## Contributing

Pull requests are encouraged. Please submit unit tests with any submitted patch. Also, please make it clear what your patch does in branch and commit messages.

1. Fork the repository
2. Create your feature branch `git checkout -b change-to-thing`
3. Commit your changes `git commit -am 'changes to the thing'`
4. Push to the branch `git push origin change-to-thing`
5. Create new pull request
6. Be patient

