package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("* * * * * *", func() {
		fmt.Println("vim-go")
	})
	c.Start()
	time.Sleep(time.Hour)
}
