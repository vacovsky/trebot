package trivia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chat-bot/bot"
	"github.com/olekukonko/tablewriter"
)

var scoresPath = "triviaScores.json"
var scores = map[string]scoreModel{}
var activeQuestion = triviaModel{
	ID:     0,
	Answer: "nil",
}

func loadScores() {
	scoreLocal := []scoreModel{}
	file, err := os.Open(scoresPath)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&scoreLocal); err != nil {
		fmt.Println(err)
	}
	for _, user := range scoreLocal {
		tmp := scoreModel{
			Name:  user.Name,
			ID:    user.ID,
			Score: user.Score,
		}
		scores[user.ID] = tmp
	}
}

func saveScores() {
	saveModel := []scoreModel{}
	for _, u := range scores {
		saveModel = append(saveModel, u)
	}
	scoresJSON, _ := json.Marshal(saveModel)
	fmt.Println(string(scoresJSON))
	err := ioutil.WriteFile(scoresPath, scoresJSON, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func renderScores() (string, error) {
	buf := &bytes.Buffer{}
	data := [][]string{}
	for _, u := range scores {
		data = append(data, []string{u.Name, strconv.Itoa(u.Score)})
	}

	table := tablewriter.NewWriter(buf) //NewWriter(os.Stdout)
	table.SetHeader([]string{"User", "Score"})

	for _, v := range data {
		table.Append(v)
	}
	table.SetAlignment(5)
	table.Render()
	return string(buf.Bytes()), nil
}

func trivia(command *bot.Cmd) (string, error) {
	if len(command.Args) < 1 {
		return "Not enough arguments!", nil
	}
	var str string
	var err error

	switch command.Args[0] {
	case "scoreboard":
		str, err = renderScores()
	case "score":
		str, err = strconv.Itoa(scores[command.User.ID].Score), nil
	case "answer":
		s := strings.Join(command.Args[1:], " ")
		str, err = checkAnswer(s, command)
	case "new":
		oldAnswer := activeQuestion.Answer
		activeQuestion, err = getTriviaClue()
		activeQuestion.ExpiresAt = time.Now().Add(time.Minute * 5)
		return fmt.Sprintf(`
---------------------------------------------------
*Previous Answer:* %s
---------------------------------------------------

===================================================
*New Question (%d) (%s):* %s
===================================================
`, oldAnswer, activeQuestion.Value, activeQuestion.Category.Title, activeQuestion.Question), err
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

func checkAnswer(answer string, command *bot.Cmd) (string, error) {
	old := activeQuestion
	if deepCheckAnswer(answer, activeQuestion.Answer) {
		activeQuestion, _ = getTriviaClue()
		activeQuestion.ExpiresAt = time.Now().Add(time.Minute * 5)
		tmp := scores[command.User.ID]
		tmp.Score += old.Value
		tmp.ID = command.User.ID
		tmp.Name = command.User.Nick
		scores[command.User.ID] = tmp
		saveScores()
		return fmt.Sprintf(`
---------------------------------------------------
%s *is correct!* ---  %s (%d)
---------------------------------------------------

===================================================
*New Question:* %s
===================================================
		`, old.Answer, command.User.Nick,
			scores[command.User.ID].Score,
			activeQuestion.Value,
			activeQuestion.Category.Title,
			activeQuestion.Question), nil

	}
	return "Try again...", nil
}

func init() {
	loadScores()

	fmt.Println("Registering Trivia...")
	bot.RegisterCommand(
		"trivia",
		"Displays a trivia question.",
		`answer {your answer}
		!trivia new
		!trivia score
		!trivia scoreboard
		`,
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
