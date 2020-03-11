package email

import (
	"crypto/tls"
	"io/ioutil"
	"net/mail"

	"github.com/asymptoter/practice-backend/base/config"
	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/sirupsen/logrus"

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

func Send(context ctx.CTX, email, message string) error {
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
	context.WithFields(logrus.Fields{
		"smtpHost":        smtpHost,
		"smtpPort":        smtpPort,
		"officialAccount": officialAccount,
		"receiver":        email,
	}).Info("Dial and send email")
	return d.DialAndSend(m)
}

func Receive(context ctx.CTX, account, password string) (string, error) {
	context.Info("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.gmail.com:993", &tls.Config{InsecureSkipVerify: false})
	if err != nil {
		context.Fatal("client.DialTLS failed ", err)
	}
	context.Info("Connected")

	// Login
	if err := c.Login(account, password); err != nil {
		context.Fatal(account+" "+password+" ", err)
	}
	defer c.Logout()
	context.Info("Logged in")

	// Select INBOX
	_, err = c.Select("INBOX", false)
	if err != nil {
		context.Fatal(err)
	}
	context.Info("Select inbox")
	//log.Println("Flags for INBOX:", mbox.Flags)

	seqset := new(imap.SeqSet)
	seqset.AddRange(0, 0)

	messages := make(chan *imap.Message, 1)
	if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchRFC822}, messages); err != nil {
		context.Fatal("Fetch failed")
	}

	res := ""
	for msg := range messages {
		for _, v := range msg.Body {
			m, err := mail.ReadMessage(v)
			if err != nil {
				context.Fatal("mail.ReadMessage ", err)
			}
			body, err := ioutil.ReadAll(m.Body)
			if err != nil {
				context.Fatal("ioutil.ReadAll ", err)
			}
			res = string(body)
		}
	}

	context.Info("Done!")
	return res, nil
}
