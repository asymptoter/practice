package email

import (
	"crypto/tls"
	"io/ioutil"
	"net/mail"

	"github.com/asymptoter/practice-backend/base/config"
	"github.com/asymptoter/practice-backend/base/ctx"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"gopkg.in/gomail.v2"
)

var (
	officialAccount  = ""
	officialPassword = ""
	smtpHost         = ""
	smtpPort         = 587
)

func Send(ctx ctx.CTX, email, message string) error {
	cfg := config.Value.Server.Email
	smtpHost = cfg.SmtpHost
	smtpPort = cfg.Port
	officialAccount = cfg.Account
	officialPassword = cfg.Password

	m := gomail.NewMessage()
	m.SetHeader("From", officialAccount)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Active practice account")
	m.SetBody("text/html", message)

	d := gomail.NewDialer(smtpHost, smtpPort, officialAccount, officialPassword)
	// TODO solve the secure issue
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	ctx.With(
		"smtpHost", smtpHost,
		"smtpPort", smtpPort,
		"officialAccount", officialAccount,
		"receiver", email,
	).Info("Dial and send email")
	return d.DialAndSend(m)
}

func Receive(ctx ctx.CTX, account, password string) (string, error) {
	ctx.Info("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.gmail.com:993", &tls.Config{InsecureSkipVerify: false})
	if err != nil {
		ctx.Fatal(err)
	}
	ctx.Info("Connected")

	// Login
	if err := c.Login(account, password); err != nil {
		ctx.Fatal(account+" "+password+" ", err)
	}
	defer c.Logout()
	ctx.Info("Logged in")

	// Select INBOX
	_, err = c.Select("INBOX", false)
	if err != nil {
		ctx.Fatal(err)
	}
	ctx.Info("Select inbox")
	//log.Println("Flags for INBOX:", mbox.Flags)

	seqset := new(imap.SeqSet)
	seqset.AddRange(0, 0)

	messages := make(chan *imap.Message, 1)
	if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchRFC822}, messages); err != nil {
		ctx.Fatal(err)
	}

	res := ""
	for msg := range messages {
		for _, v := range msg.Body {
			m, err := mail.ReadMessage(v)
			if err != nil {
				ctx.Fatal(err)
			}
			body, err := ioutil.ReadAll(m.Body)
			if err != nil {
				ctx.Fatal(err)
			}
			res = string(body)
		}
	}

	ctx.Info("Done!")
	return res, nil
}
