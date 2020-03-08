package main

import (
	"log"
	"net/smtp"

	"github.com/asymptoter/practice-backend/base/config"
)

var (
	from = config.Value.Server.Email.Account
	pass = config.Value.Server.Email.Password
)

func main() {
	to := "asymptotion@gmail.com"
	msg := "test"
	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")
	if err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg)); err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}
