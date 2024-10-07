package utils

import (
	"bytes"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/go-gomail/gomail"
)

func SendEmail(htmlTemplate, targetEmail, subjectMessage, attachmentPath string, fill any) {
	host := config.LoadConfig().SmtpHost
	username := config.LoadConfig().SmtpUsername
	password := config.LoadConfig().SmtpPassword
	port := config.LoadConfig().SmtpPort

	htmlFile, err := os.ReadFile("assets/html/" + htmlTemplate + ".html")
	if err != nil {
		log.Println("Failed to read file:", err.Error())
		return
	} else {
		log.Println("Success to read file")
	}

	tmpl, err := template.New("emailTemplate").Parse(string(htmlFile))
	if err != nil {
		log.Println("Failed to parse template:", err.Error())
		return
	} else {
		log.Println("Success to parse template")
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, fill)
	if err != nil {
		log.Println("Failed to fill in template:", err.Error())
		return
	} else {
		log.Println("Success to fill in template")
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", username)
	mailer.SetHeader("To", strings.ToLower(targetEmail))
	mailer.SetHeader("Subject", subjectMessage)
	mailer.SetBody("text/html", tpl.String())
	if attachmentPath != "" {
		mailer.Attach(attachmentPath)
	}

	d := gomail.NewDialer(host, port, username, password)
	if err := d.DialAndSend(mailer); err != nil {
		log.Println("Failed to send email. Error: ", err.Error())
		return
	} else {
		log.Println("Success to send email. Template: "+htmlTemplate+".html. Target Email: ", strings.ToLower(targetEmail))
	}
}
