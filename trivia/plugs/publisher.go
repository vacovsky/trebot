package plugs

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

// This is for use with neatlights

var CorrectAnswerPS = `{'direction': 2, 'iterations': 1, 'offset': 1, 'senselight': 0, 'speed': 0.01, 'css3_colors': ['red', 'yellow', 'blue', 'orange', 'purple', 'green', 'teal', 'fuchsia', 'white'], 'brightness': 10, 'reverse_after': 500, 'method_name': 'party', 'cleanup': 1, 'style_name': 'Diamonds', 'led_count': 150}`

var IncorrectAnswerPS = `{'senselight': 0, 'brightness': 10, 'color': ['firebrick'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 150}`

var NewQuestionPS = `{'senselight': 0, 'brightness': 10, 'color': ['skyblue'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 150}`

// Publish sends a blob to a redis pubsub channel
func Publish(channel, message string) {
	channel = `ws2812b_0`
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("TREBOT_REDIS_ADDR"),
		Password: os.Getenv("TREBOT_REDIS_SECRET"),
	})
	defer client.Close()
	err := client.Publish(channel, message).Err()
	if err != nil {
		log.Println(err)
	}
}
