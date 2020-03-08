package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/mail"

	"github.com/asymptoter/practice-backend/base/config"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func main() {
	config.Init()
	cfg := config.Value
	account := cfg.Server.Testing.Email.Account[:8]
	password := cfg.Server.Testing.Email.Password
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.gmail.com:993", &tls.Config{InsecureSkipVerify: false})
	if err != nil {
		log.Fatal("client.DialTLS failed ", err)
	}
	log.Println("Connected")

	// Login
	if err := c.Login(account, password); err != nil {
		log.Fatal(account+" "+password+" ", err)
	}
	defer c.Logout()
	log.Println("Logged in")

	// Select INBOX
	_, err = c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Select inbox")
	//log.Println("Flags for INBOX:", mbox.Flags)

	seqset := new(imap.SeqSet)
	seqset.AddRange(0, 0)

	messages := make(chan *imap.Message, 1)
	// c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	//if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchRFC822Text}, messages); err != nil {
	if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchRFC822Text}, messages); err != nil {
		log.Fatal("Fetch failed")
	}

	log.Println("Last message:")
	for msg := range messages {
		//log.Println("* " + msg.Envelope.Subject)
		log.Println("range body")
		for name, v := range msg.Body {
			//log.Println("@", name, v)
			r := msg.GetBody(name)
			log.Println("#", r)
			m, err := mail.ReadMessage(v)
			if err != nil {
				log.Fatal("mail.ReadMessage ", err)
			}
			body, err := ioutil.ReadAll(m.Body)
			if err != nil {
				log.Fatal("ioutil.ReadAll ", err)
			}
			log.Println(body)
		}
	}

	log.Println("Done!")
}
