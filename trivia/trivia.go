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
	// fmt.Println(string(scoresJSON))
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
	// fmt.Println(string(buf.Bytes()))
	return "```" + string(buf.Bytes()) + "```", nil
}
func showAbout() (string, error) {
	return `
	> This plugin for go-chat-bot (https://github.com/go-chat-bot/bot) leverages jService to provide every Jeopardy question ever.  Thanks to the person who made that! 
	> By Joe Vacovsky Jr.
	> Bot source code located at https://github.com/vacoj/trebot
	> Submit bugs/issues at https://github.com/vacoj/trebot/issues
	`, nil
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
	case "about":
		str, err = showAbout()
	case "stats":
		str, err = showStats(command)
		str = "```" + str + "```"
	case "score":
		str, err = strconv.Itoa(scores[command.User.ID].Score), nil
		str = command.User.Nick + ": " + str
	case "answer":
		s := strings.Join(command.Args[1:], " ")
		str, err = checkAnswer(s, command)
	case "new":
		oldAnswer := activeQuestion[command.Channel].Answer
		q, err := getTriviaClue()
		q.ExpiresAt = time.Now().Add(time.Minute * 5)
		activeQuestion[command.Channel] = q
		return fmt.Sprintf(`
:unamused:  *Previous Answer:* 
> %s

:question:  *New Question (%s for %d):* 
> %s
`,
			oldAnswer,
			activeQuestion[command.Channel].Category.Title,
			activeQuestion[command.Channel].Value,
			activeQuestion[command.Channel].Question), err
	default:
		return fmt.Sprintf(`
:question:  *Current Question (%s for %d):* 
> %s
		`,
			activeQuestion[command.Channel].Category.Title,
			activeQuestion[command.Channel].Value,
			activeQuestion[command.Channel].Question), nil
	}

	if err != nil {
		return fmt.Sprintf("Error: %s", err), nil
	}
	return str, nil
}

func scrubStrings(input string) string {
	cleanups := []string{
		`^(an )`,
		`^(the )`,
		`(<.>)`,
		`^(a )`,
		`(<..>)`,
		`(\(|\))`,
		`(\")`,
		`(\\')`,
	}
	byteAnswer := []byte(input)

	for _, c := range cleanups {
		// fmt.Println(string(byteAnswer))
		rex := regexp.MustCompile(c)
		byteAnswer = rex.ReplaceAll(byteAnswer, []byte(""))
	}
	// fmt.Println(string(byteAnswer))

	return strings.ToLower(string(byteAnswer))
}

func showStats(cmd *bot.Cmd) (string, error) {
	prettyScoreModel := fmt.Sprintf(`
	Player Name: %s
	Total Score: %d
	Total Correct Answers: %d
	Total Wrong Answers: %d
	Total Requested New Questions: %d
`,
		scores[cmd.User.ID].Name,
		scores[cmd.User.ID].Score,
		scores[cmd.User.ID].CorrectAnswers,
		scores[cmd.User.ID].WrongAnswers,
		scores[cmd.User.ID].NewQuestionRequests,
	)
	return prettyScoreModel, nil
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
	q[0].Answer = scrubStrings(q[0].Answer)

	// fmt.Println(q[0].Question, " ***** ", q[0].Answer)
	return q[0], nil
}

func checkAnswer(answer string, command *bot.Cmd) (string, error) {
	old := activeQuestion[command.Channel]
	if deepCheckAnswer(answer, activeQuestion[command.Channel].Answer) {
		q, _ := getTriviaClue()
		q.ExpiresAt = time.Now().Add(time.Minute * 5)
		activeQuestion[command.Channel] = q
		tmp := scores[command.User.ID]
		tmp.CorrectAnswers++
		tmp.Score += old.Value
		tmp.ID = command.User.ID
		tmp.Name = command.User.Nick
		scores[command.User.ID] = tmp
		saveScores()
		return fmt.Sprintf(`
:moneybag:  *%s* is correct! ---  %s (%d)

:question:  *New Question (%s for %d):* 
> %s
		`, old.Answer,
			command.User.Nick,
			scores[command.User.ID].Score,
			activeQuestion[command.Channel].Category.Title,
			activeQuestion[command.Channel].Value,
			activeQuestion[command.Channel].Question), nil
	}
	tmp := scores[command.User.ID]
	tmp.WrongAnswers++
	scores[command.User.ID] = tmp
	saveScores()
	return "Try again...", nil
}

func init() {
	loadScores()

	fmt.Println("Registering Trivia...")
	bot.RegisterCommand(
		"trivia",
		"Displays a trivia question.",
		`answer {your answer}
		!trivia new (stops current question, and pitches a new question)
		!trivia score (shows your score)
		!trivia scoreboard (shows all players' scores, ranked from highest -> lowest)
		!trivia stats (show's your stats)
		!trivia about (shows information related to this trivia bot)
		`,
		trivia)
}

func deepCheckAnswer(providedAnswer, realAnswer string) bool {
	lowerP, lowerR := strings.ToLower(providedAnswer), strings.ToLower(realAnswer)
	// fmt.Print1ln(lowerP, ":", lowerR, ":", string(byteAnswer))

	if len([]byte(lowerP)) >= 5 && strings.Contains(lowerR, lowerP) {
		return true
	} else if realAnswer == providedAnswer {
		return true
	}
	return false
}
