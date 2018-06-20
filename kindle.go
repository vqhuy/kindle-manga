package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/c633/kindle-manga/util"
	gomail "gopkg.in/gomail.v2"
)

const subject = "Kindle-Manga"

type Mail struct {
	From     string
	To       string
	Subject  string
	Password string
}

func getMailSettings() *Mail {
	var msg [3]string
	var inputs [3]string

	msg[0] = "Enter your approved Send-To-Kindle email"
	msg[1] = "Enter your email password (or your app password if you've enabled 2FA)"
	msg[2] = "Enter your Kindle's email"

	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < 3; i++ {
		fmt.Println(msg[i])
		fmt.Print("-> ")
		inputs[i], _ = reader.ReadString('\n')
		// convert CRLF to LF
		inputs[i] = strings.Replace(inputs[i], "\n", "", -1)
	}

	return &Mail{
		From:     inputs[0],
		Password: inputs[1],
		To:       inputs[2],
		Subject:  subject,
	}
}

func restoreMailSettings() (*Mail, error) {
	mailSettings := &Mail{}
	mailFile := filepath.Join(configDir, "mail.json")
	err := util.LoadJSONFromFile(mailFile, mailSettings)
	if err != nil {
		log.Println(err)
		mailSettings = getMailSettings()

		err = util.SaveJSONToFile(mailFile, mailSettings)
		if err != nil {
			return nil, err
		}
	}
	return mailSettings, nil
}

func sendToKindle(mail *Mail, file string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mail.From)
	m.SetHeader("To", mail.To)
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/plain", "")
	m.Attach(file)

	d := gomail.NewPlainDialer("smtp.gmail.com", 587, mail.From, mail.Password)
	return d.DialAndSend(m)
}
