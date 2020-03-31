package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type TriviaMode int

const (
	TriviaModePlayAll   TriviaMode = iota // Count correct quizzes
	TriviaModeNoWrong                     // Wrong answer then game end
	TriviaModeTimeCount                   // Count correct quizzes in time limit
)

type Quiz struct {
	ID       int64          `json:"ID" db:"id"`
	Content  string         `json:"content" db:"content"`
	ImageURL string         `json:"imageURL" db:"image_url"`
	Options  pq.StringArray `json:"options" db:"options"`
	Answer   string         `json:"answer" db:"answer"`
	Creator  uuid.UUID      `json:"creator" db:"creator"`
	Category string         `json:"category" db:"category"`
}

type Game struct {
	ID        uuid.UUID     `json:"ID" db:"id"`
	Name      string        `json:"name" db:"name"`
	QuizIDs   pq.Int64Array `json:"quizIDs" db:"quiz_ids"`
	Mode      TriviaMode    `json:"mode" db:"mode"`
	CountDown int           `json:"countDown" db:"count_down"`
	Creator   uuid.UUID     `json:"creator" db:"creator"`
}

type GameStatus struct {
	Name      string        `json:"name"`
	QuizNo    int           `json:"quizNo"`
	QuizIDs   pq.Int64Array `json:"quizIDs"`
	Answers   []string      `json:"answers"`
	Mode      TriviaMode    `json:"mode"`
	CountDown int           `json:"countDown"`
}

type GameResult struct {
	GameID       int64         `json:"ID" db:"id"`
	Player       uuid.UUID     `json:"player" db:"player"`
	TimeSpent    time.Duration `json:"timeSpent" db:"time_spent"`
	CorrectCount int           `json:"correctCount" db:"correct_count"`
}
