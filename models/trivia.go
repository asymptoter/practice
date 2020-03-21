package models

import (
	"time"

	"github.com/google/uuid"
)

type TriviaMode int

const (
	PlayAll   TriviaMode = iota // Count correct quizzes
	NoWrong                     // Wrong answer then game end
	TimeCount                   // Count correct quizzes in time limit
)

type Quiz struct {
	ID       int64     `json:"ID" db:"id"`
	Content  string    `json:"Content" db:"content"`
	ImageURL string    `json:"ImageURL" db:"image_url"`
	Options  []string  `json:"Options" db:"options"`
	Answer   string    `json:"Answer" db:"answer"`
	Creator  uuid.UUID `json:"Creator" db:"creator"`
}

type Game struct {
	ID        int64      `json:"ID" db:"id"`
	QuizIDs   []int64    `db:"quiz_ids"`  // Used in db
	Quizzes   []Quiz     `json:"Quizzes"` // Used in response
	Mode      TriviaMode `json:"Mode" db:"mode"`
	CountDown int        `json:"CountDown" db:"count_down"`
	Creator   uuid.UUID  `json:"Creator" db:"creator"`
}

type GameResult struct {
	GameID       int64         `json:"ID" db:"id"`
	Player       uuid.UUID     `json:"Player" db:"player"`
	TimeSpent    time.Duration `json:"TimeSpent" db:"time_spent"`
	CorrectCount int           `json:"CorrectCount" db:"correct_count"`
}
