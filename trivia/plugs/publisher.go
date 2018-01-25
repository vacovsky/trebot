package plugs

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

// This is for use with neatlights

var CorrectAnswerPS = `{'senselight': 0, 'brightness': 10, 'color': ['green'], 'method_name': 'room_lighting', 'style_name': 'Room Lighting', 'led_count': 150}`

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
