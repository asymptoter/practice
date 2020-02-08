package main

import (
	"log"
	"net/smtp"
)

func main() {
	from := "k4kugybv@gmail.com"
	pass := "Pxb43gEq"
	to := "asymptotion@gmail.com"
	msg := "test"
	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")
	if err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg)); err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}
