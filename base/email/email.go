package email

import (
	"crypto/tls"

	"github.com/asymptoter/practice-backend/base/config"
	"github.com/asymptoter/practice-backend/base/ctx"

	"gopkg.in/gomail.v2"
)

var (
	officialAccount  = ""
	officialPassword = ""
	smtpAddress      = ""
	smtpPort         = 587
)

func init() {
	smtpAddress = config.Value.Server.Email.Address
	smtpPort = config.Value.Server.Email.Port
	officialAccount = config.Value.Server.Email.OfficialAccount
	officialPassword = config.Value.Server.Email.OfficialPassword
}

func Send(context ctx.CTX, email, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", officialAccount)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Active practice account")
	m.SetBody("text/html", message)

	d := gomail.NewDialer(smtpAddress, smtpPort, officialAccount, officialPassword)
	// TODO solve the secure issue
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}
