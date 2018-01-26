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
	"sync"
	"time"

	"./admin"
	"./plugs"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-chat-bot/bot"
	"github.com/olekukonko/tablewriter"
)

var (
	scoresPath       = "triviaScores.json"
	bansPath         = "triviaBans.json"
	adminsPath       = "triviaAdmins.json"
	scores           = map[string]scoreModel{}
	bans             = map[string]admin.BanModel{}
	admins           = map[string]admin.AdminModel{}
	activeQuestion   = map[string]triviaModel{}
	previousQuestion = map[string]triviaModel{}
	scoreLock        = sync.Mutex{}
)

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
			Name:                user.Name,
			ID:                  user.ID,
			Score:               user.Score,
			CorrectAnswers:      user.CorrectAnswers,
			NewQuestionRequests: user.NewQuestionRequests,
			WrongAnswers:        user.WrongAnswers,
		}
		scores[user.ID] = tmp
	}
}

func saveScores() {
	saveModel := []scoreModel{}
	scoreLock.Lock()
	for _, u := range scores {
		saveModel = append(saveModel, u)
	}
	scoreLock.Unlock()
	scoresJSON, _ := json.MarshalIndent(saveModel, "", "    ")
	err := ioutil.WriteFile(scoresPath, scoresJSON, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func calcAccuracy(correct, incorrect int) string {
	var result float64
	div := float64(correct) + float64(incorrect)
	if div == 0.0 {
		div = 1
	}
	result = float64(correct) / div
	// fmt.Println(result)
	return strconv.FormatFloat(result, 'f', 3, 64)
}

func renderScores() (string, error) {
	buf := &bytes.Buffer{}
	data := [][]string{}
	scoreLock.Lock()
	for _, u := range scores {
		thisUser := []string{
			u.Name,
			strconv.Itoa(u.Score),
			strconv.Itoa(u.CorrectAnswers),
			strconv.Itoa(u.WrongAnswers),
			strconv.Itoa(u.NewQuestionRequests),
			calcAccuracy(u.CorrectAnswers, u.WrongAnswers),
		}
		data = append(data, thisUser)
	}
	scoreLock.Unlock()

	table := tablewriter.NewWriter(buf) //NewWriter(os.Stdout)
	table.SetHeader([]string{"Rank", "User", "Score", "Correct", "Wrong", "New", "Accuracy"})
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
*About Trebot*

Commands:

> *!trivia answer {your answer}* Alternatively, you can simply type in your answer without the leading *!trivia answer* portion
> *!trivia new* (stops current question, and pitches a new question)

> *!trivia scoreboard* (shows all players' scores and stats, ranked from highest -> lowest scores)
> *!trivia stats* (shows your stats)

> *!trivia about* (shows information related to this trivia bot)
> *!trivia rules* (shows information about how answers are matched, and about scoring rules)

General Info:

> This plugin for go-chat-bot (https://github.com/go-chat-bot/bot) leverages jService (http://jservice.io) to provide every Jeopardy question ever.  Thanks to the person who made that! 
> by Joe Vacovsky Jr. (https://github.com/vacoj)
> Bot source code located at https://github.com/vacoj/trebot
> Submit bugs/issues at https://github.com/vacoj/trebot/issues
`, nil
}

func showRules() (string, error) {
	return `
*Trebot Rules*

Matching:

> Actual answers have the following characters / words removed before comparison (players' answers are not modified)
>     '  "  .  (  )  an  the  a 
> Partial matches are acceptable, as long as the partial answer given is 5 or more characters
> Players' answers with less than 5 characters must be an *exact* match to the actual answer
> If you want the matching rules to be better, feel free to contribute (https://github.com/vacoj/trebot)

Scoring:

> All text not prefaced by !trivia will be considered an answer
> Wrong answers are penalized as a percentage of the question's value to discourage volume-guessing
> Bots are not allowed to play
`, nil
}

func answer(c *bot.PassiveCmd) (string, error) {
	if c.User.IsBot {
		return fmt.Sprintf("Sorry %s, bots are not allowed to play.", c.User.Nick), nil
	}

	for b := range bans {
		if bans[b].ID == c.User.ID {
			return fmt.Sprintf(":no_entry_sign: Sorry, %s, you're not allowed to play.  If you think this is unfair, please check with one of the admins.", c.User.Nick), nil
		}
	}
	return checkAnswerSilently(c.Raw, c)
}

func checkIfAdmin(userID string) bool {
	for u := range admins {
		if admins[u].ID == userID {
			return true
		}
	}
	return false
}

func banUser(userID string, command *bot.Cmd) {
	bans[userID] = admin.BanModel{
		ID:           userID,
		BannedSince:  time.Now(),
		BannedByID:   command.User.ID,
		BannedByName: command.User.Nick,
	}
	admin.SaveBans(bansPath, bans)
}

func unbanUser(userID string) {
	tmpBans := map[string]admin.BanModel{}

	for b := range bans {
		if bans[b].ID != userID {
			tmpBans[bans[b].ID] = bans[b]
		}
	}
	bans = tmpBans
	admin.SaveBans(bansPath, tmpBans)
}

func trivia(command *bot.Cmd) (string, error) {
	if len(command.Args) < 1 {
		return "Not enough arguments!", nil
	}

	for b := range bans {
		if bans[b].ID == command.User.ID {
			return fmt.Sprintf(":no_entry_sign: Sorry, %s, you're not allowed to play.  If you think this is unfair, please check with one of the admins.", command.User.Nick), nil
		}
	}

	var str string
	var err error

	switch command.Args[0] {

	case "ban":
		if len(command.Args) == 2 {
			if checkIfAdmin(command.User.ID) {
				for b := range bans {
					if bans[b].ID == command.Args[1] {
						return fmt.Sprintf(":jomoji-dealwithit-animated: %s is already banned.", command.Args[1]), nil
					}
				}
				banUser(command.Args[1], command)
				return fmt.Sprintf(":jomoji-devil: %s has been banned.  !trivia unban %s to unban.", command.Args[1], command.Args[1]), nil
			}
			return "Only administrators can ban and unban.", nil
		}

	case "unban":
		if len(command.Args) == 2 {
			if checkIfAdmin(command.User.ID) {
				unbanUser(command.Args[1])
				return fmt.Sprintf(":jomoji-angel: %s has been unbanned.", command.Args[1]), nil
			}
			return "Only administrators can ban and unban.", nil
		}

	case "scoreboard":
		str, err = renderScores()

	case "about":
		str, err = showAbout()

	case "rules":
		str, err = showRules()

	case "stats":

		str, err = showStats(command)
		str = "```" + str + "```"

	case "answer":
		if command.User.IsBot {
			return fmt.Sprintf("Sorry %s, bots are not allowed to play.", command.User.Nick), nil
		}

		s := strings.Join(command.Args[1:], " ")
		str, err = checkAnswer(s, command)

	case "new":
		oldAnswer := activeQuestion[command.Channel].Answer
		q, err := getTriviaClue()
		q.ExpiresAt = time.Now().Add(time.Minute * 5)
		activeQuestion[command.Channel] = q
		tmp := scores[command.User.ID]
		tmp.NewQuestionRequests++
		scores[command.User.ID] = tmp
		saveScores()
		plugs.Publish(``, plugs.NewQuestionPS)
		return fmt.Sprintf(`
:unamused:  Previous Answer: 
> *%s*

:question:  New Question ([%d] *%s* for *%d*): 
> *%s*
`,
			oldAnswer,
			activeQuestion[command.Channel].Airdate.Year(),
			activeQuestion[command.Channel].Category.Title,
			activeQuestion[command.Channel].Value,
			activeQuestion[command.Channel].Question), err

	default:
		return showAbout()
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
Total New Question Requests: %d
Accuracy: %s
`,
		scores[cmd.User.ID].Name,
		scores[cmd.User.ID].Score,
		scores[cmd.User.ID].CorrectAnswers,
		scores[cmd.User.ID].WrongAnswers,
		scores[cmd.User.ID].NewQuestionRequests,
		calcAccuracy(scores[cmd.User.ID].CorrectAnswers, scores[cmd.User.ID].WrongAnswers),
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

:question:  New Question ([%d] *%s* for *%d*): 
> *%s*
		`, old.Answer,
			command.User.Nick,
			scores[command.User.ID].Score,
			activeQuestion[command.Channel].Airdate.Year(),
			activeQuestion[command.Channel].Category.Title,
			activeQuestion[command.Channel].Value,
			activeQuestion[command.Channel].Question), nil
	}
	tmp := scores[command.User.ID]
	tmp.WrongAnswers++
	tmp.Score -= old.Value / 10
	scores[command.User.ID] = tmp
	saveScores()
	return "Try again...", nil
}

func checkAnswerSilently(answer string, command *bot.PassiveCmd) (string, error) {
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

:question:  New Question ([%d] *%s* for *%d*): 
> *%s*
		`, old.Answer,
			command.User.Nick,
			scores[command.User.ID].Score,
			activeQuestion[command.Channel].Airdate.Year(),
			activeQuestion[command.Channel].Category.Title,
			activeQuestion[command.Channel].Value,
			activeQuestion[command.Channel].Question), nil
	}
	tmp := scores[command.User.ID]
	tmp.WrongAnswers++
	tmp.Score -= old.Value / 10
	scores[command.User.ID] = tmp
	saveScores()
	return "", nil
}

func init() {
	loadScores()
	admins = admin.LoadAdmins(adminsPath)
	bans = admin.LoadBans(bansPath)

	spew.Dump(bans)
	spew.Dump(admins)

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
	bot.RegisterPassiveCommand(``, answer)
}

func deepCheckAnswer(providedAnswer, realAnswer string) bool {
	lowerP, lowerR := strings.ToLower(strings.Trim(providedAnswer, " ")), strings.ToLower(strings.Trim(realAnswer, " "))
	if len([]byte(lowerP)) >= 5 && strings.Contains(lowerR, lowerP) {
		plugs.Publish(``, plugs.CorrectAnswerPS)
		return true
	} else if lowerR == lowerP {
		plugs.Publish(``, plugs.CorrectAnswerPS)
		return true
	}
	plugs.Publish(``, plugs.IncorrectAnswerPS)
	return false
}
