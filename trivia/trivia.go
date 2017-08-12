package trivia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chat-bot/bot"
	"github.com/olekukonko/tablewriter"
)

var scoresPath = "triviaScores.json"
var scores = map[string]scoreModel{}
var activeQuestion = map[string]triviaModel{}
var previousQuestion = map[string]triviaModel{}

func loadScores() {
	scoreLocal := []scoreModel{}
	file, err := os.Open(scoresPath)
	defer file.Close()
	tried := false
	if err != nil && !tried {
		tried = true
		saveScores()
		loadScores()
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
	table.SetHeader([]string{"Rank", "User", "Score"})
	sort.Slice(data, func(i, j int) bool {
		d1, _ := strconv.Atoi(data[i][1])
		d2, _ := strconv.Atoi(data[j][1])
		return d1 > d2
	})

	for i, v := range data {
		v = append([]string{strconv.Itoa(i + 1)}, v...)
		table.Append(v)
	}
	table.Render()
	fmt.Println(string(buf.Bytes()))
	return "```" + string(buf.Bytes()) + "```", nil
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
		oldAnswer := previousQuestion[command.Channel].Answer
		q, err := getTriviaClue()
		q.ExpiresAt = time.Now().Add(time.Minute * 5)
		activeQuestion[command.Channel] = q
		return fmt.Sprintf(`
---------------------------------------------------
*Previous Answer:* %s
---------------------------------------------------

===================================================
*New Question (%s for %d):* %s
===================================================
`,
			oldAnswer,
			activeQuestion[command.Channel].Category.Title,
			activeQuestion[command.Channel].Value,
			activeQuestion[command.Channel].Question), err
	default:
		return "Not enough arguments.", nil
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
	if q[0].Value == 0 {
		q[0].Value = 5000
	}
	fmt.Println(q[0].Question, " |||| ", q[0].Answer)
	return q[0], nil
}

func checkAnswer(answer string, command *bot.Cmd) (string, error) {
	old := activeQuestion[command.Channel]
	if deepCheckAnswer(answer, activeQuestion[command.Channel].Answer) {
		q, _ := getTriviaClue()
		q.ExpiresAt = time.Now().Add(time.Minute * 5)
		activeQuestion[command.Channel] = q
		tmp := scores[command.User.ID]
		tmp.Score += old.Value
		tmp.ID = command.User.ID
		tmp.Name = command.User.Nick
		scores[command.User.ID] = tmp
		saveScores()
		return fmt.Sprintf(`
---------------------------------------------------
*%s* is correct! ---  %s (%d)
---------------------------------------------------

===================================================
*New Question (%s for %d):* %s
===================================================
		`, old.Answer,
			command.User.Nick,
			scores[command.User.ID].Score,
			activeQuestion[command.Channel].Category.Title,
			activeQuestion[command.Channel].Value,
			activeQuestion[command.Channel].Question), nil

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
	// fmt.Print1ln(lowerP, ":", lowerR, ":", string(byteAnswer))

	if len([]byte(lowerP)) >= 5 && strings.Contains(lowerR, lowerP) {
		return true
	} else if realAnswer == providedAnswer {
		return true
	}
	return false
}
