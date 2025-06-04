package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/powiedl/myGoWebApplication/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	go func() {
		for {
			msg := <- app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
  log.Printf("Send Mail from:%s to %s with subject '%s'\n",m.From,m.To,m.Subject)
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client,err := server.Connect()
	if err != nil {
		errorLog.Println("Error connecting to the mail server:",err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {	
		email.SetBody(mail.TextHTML,m.Content)
	} else {
		data,err := os.ReadFile(fmt.Sprintf("%sstatic/email/templates/%s",app.Basedir,m.Template))
		if err != nil {
			app.ErrorLog.Printf("Error occured while reading template '%s':%v",m.Template,err)
		}
		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate,"[%E_MAIL_CONTENT%]",m.Content,1)
		email.SetBody(mail.TextHTML,msgToSend)
	}
	err = email.Send(client)
	if err != nil {
		log.Printf("Error sending mail from:%s to %s with subject '%s':%v\n",m.From,m.To,m.Subject,err)
	} else {
		log.Println("E-Mail successfully transmitted to the mail server")
	}
}