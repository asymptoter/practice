package models

import "time"

type TriviaMode int

const (
	PlayAll TriviaMode = iota
	NoWrong
	TimeCount
)

type Quiz struct {
	ID      int64
	Content string
	Options []string
	Answer  int
	Creator string
}

type Game struct {
	ID        int64
	Quizzes   []Quiz
	Mode      TriviaMode
	CountDown int
}

type GameResult struct {
	GameID       int64
	TimeSpent    time.Duration
	CorrectCount int
}
