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
	Content  string    `json:"content" db:"content"`
	ImageURL string    `json:"imageURL" db:"image_url"`
	Options  []string  `json:"options" db:"options"`
	Answer   string    `json:"answer" db:"answer"`
	Creator  uuid.UUID `json:"creator" db:"creator"`
	Category string    `json:"category" db:"category"`
}

type Game struct {
	ID        int64      `json:"ID" db:"id"`
	QuizIDs   []int64    `db:"quiz_ids"`  // Used in db
	Quizzes   []Quiz     `json:"quizzes"` // Used in response
	Mode      TriviaMode `json:"mode" db:"mode"`
	CountDown int        `json:"countDown" db:"count_down"`
	Creator   uuid.UUID  `json:"creator" db:"creator"`
}

type GameResult struct {
	GameID       int64         `json:"ID" db:"id"`
	Player       uuid.UUID     `json:"player" db:"player"`
	TimeSpent    time.Duration `json:"timeSpent" db:"time_spent"`
	CorrectCount int           `json:"correctCount" db:"correct_count"`
}
