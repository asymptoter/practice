package email

import (
	"github.com/asymptoter/geochallenge-backend/base/ctx"

	"gopkg.in/gomail.v2"
)

var (
	officialEmailAccount  = ""
	officialEmailPassword = ""
)

func InitialEmailSetting(account, password string) {
	officialEmailAccount = account
	officialEmailPassword = password
}

func Send(context ctx.CTX, email, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", officialEmailAccount)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Active geochallenge account")
	m.SetBody("text/html", message)

	d := gomail.NewDialer("smtp.google.com", 587, officialEmailAccount, officialEmailPassword)

	return d.DialAndSend(m)
}
