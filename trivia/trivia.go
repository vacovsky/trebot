package trivia

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-chat-bot/bot"
)

var activeQuestion = triviaModel{
	ID:     0,
	Answer: "nil",
}

func trivia(command *bot.Cmd) (string, error) {
	if len(command.Args) < 1 {
		return "Not enough arguments!", nil
	}

	var str string
	var err error

	switch command.Args[0] {
	case "answer":
		s := strings.Join(command.Args[1:], " ")
		str, err = checkAnswer(s)
	case "new":
		oldAnswer := activeQuestion.Answer
		activeQuestion, err = getTriviaClue()
		activeQuestion.ExpiresAt = time.Now().Add(time.Minute * 5)
		fmt.Println(activeQuestion.Question, " |||| ", activeQuestion.Answer)
		return fmt.Sprintf(`
---------------------
Previous Answer: %s
---------------------

+++++++++++++++++++++
New Question: %s
+++++++++++++++++++++
`, oldAnswer, activeQuestion.Question), err
	default:
		// if activeQuestion
		if activeQuestion.ID == 0 || time.Now().Unix() > activeQuestion.ExpiresAt.Unix() {
			activeQuestion, err = getTriviaClue()
			activeQuestion.ExpiresAt = time.Now().Add(time.Minute * 5)
			return activeQuestion.Question, err
		} else {
			return fmt.Sprintf(`
Current active question is: 
- %s

- Expires at: %s
`, activeQuestion.Question, activeQuestion.ExpiresAt), nil

		}
	}
	if err != nil {
		return fmt.Sprintf("Error: %s", err), nil
	}
	return str, nil
}

func getTriviaClue() (triviaModel, error) {
	jservice := "http://jservice.io/api/random"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", jservice, nil)
	r, _ := client.Do(req)
	body, _ := ioutil.ReadAll(r.Body)

	q := []triviaModel{}
	_ = json.Unmarshal(body, &q)
	return q[0], nil
}

func checkAnswer(answer string) (string, error) {
	if strings.ToLower(answer) == strings.ToLower(activeQuestion.Answer) {
		activeQuestion, _ = getTriviaClue()
		activeQuestion.ExpiresAt = time.Now().Add(time.Minute * 5)
		return fmt.Sprintf(`
---------------------
%s is correct!
---------------------

+++++++++++++++++++++		
New Question: %s
+++++++++++++++++++++
		`, answer, activeQuestion.Question), nil
		fmt.Println(activeQuestion.Question, " |||| ", activeQuestion.Answer)

	}
	return "Try again...", nil
}

func init() {
	fmt.Println("Registering Trivia...")
	bot.RegisterCommand(
		"trivia",
		"Displays a trivia question.",
		"answer {your answer}",
		trivia)
}
