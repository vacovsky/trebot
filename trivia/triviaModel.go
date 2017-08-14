package trivia

import "time"

type triviaModel struct {
	ID           int         `json:"id"`
	ExpiresAt    time.Time   `json:"expires_at"`
	Answer       string      `json:"answer"`
	Question     string      `json:"question"`
	Value        int         `json:"value"`
	Airdate      time.Time   `json:"airdate"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	CategoryID   int         `json:"category_id"`
	GameID       interface{} `json:"game_id"`
	InvalidCount interface{} `json:"invalid_count"`
	Category     struct {
		ID         int       `json:"id"`
		Title      string    `json:"title"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		CluesCount int       `json:"clues_count"`
	} `json:"category"`
}

type scoreModel struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Score               int    `json:"score"`
	CorrectAnswers      int    `json:"correct_answers"`
	WrongAnswers        int    `json:"wrong_answers"`
	NewQuestionRequests int    `json:"new_question_requests"`
}
