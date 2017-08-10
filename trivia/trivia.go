package trivia

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
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
	fmt.Println(q[0].Question, " |||| ", q[0].Answer)
	return q[0], nil
}

func checkAnswer(answer string) (string, error) {
	oldAnswer := activeQuestion.Answer
	if deepCheckAnswer(answer, activeQuestion.Answer) {
		activeQuestion, _ = getTriviaClue()
		activeQuestion.ExpiresAt = time.Now().Add(time.Minute * 5)
		return fmt.Sprintf(`
---------------------
%s is correct!
---------------------

+++++++++++++++++++++		
New Question: %s
+++++++++++++++++++++
		`, oldAnswer, activeQuestion.Question), nil
	}
	return "Try again...", nil
}

func init() {
	fmt.Println("Registering Trivia...")
	bot.RegisterCommand(
		"trivia",
		"Displays a trivia question.",
		`answer {your answer}
		!trivia new`,
		trivia)
}

func deepCheckAnswer(providedAnswer, realAnswer string) bool {
	cleanups := []string{
		`^(an )/g`,
		`^(the )/g`,
		`(<.>)|(<..>)/g`,
		`(.*)/g`,
		`^(a )/g`,
		`(<..>)/g`,
		`(")/g`,
	}
	byteAnswer := []byte(realAnswer)

	for _, c := range cleanups {
		fmt.Println()
		rex := regexp.MustCompile(c)
		byteAnswer = rex.ReplaceAll(byteAnswer, []byte(""))
	}
	lowerP, lowerR := strings.ToLower(providedAnswer), strings.ToLower(realAnswer)
	fmt.Println(lowerP, ":", lowerR, ":", string(byteAnswer))

	if len([]byte(lowerP)) >= 5 && strings.Contains(lowerR, lowerP) {
		return true
	} else if realAnswer == providedAnswer {
		return true
	}
	return false
}
