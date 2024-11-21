package smtp

import (
	"bytes"
	"log"
	"strings"
	"text/template"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/pkg/upload"
	"github.com/go-gomail/gomail"
)

func SendEmail(htmlTemplate, senderEmail, targetEmail, subjectMessage, attachmentPath string, fill any) {
	host := config.LoadConfig().SmtpHost
	username := config.LoadConfig().SmtpUsername
	password := config.LoadConfig().SmtpPassword
	port := config.LoadConfig().SmtpPort

	if senderEmail == "" {
		senderEmail = username
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", senderEmail)
	mailer.SetHeader("To", strings.ToLower(targetEmail))
	mailer.SetHeader("Subject", subjectMessage)

	htmlFile, _, err := upload.GetRemoteFile("/assets/html/" + htmlTemplate + ".html")
	if err != nil {
		log.Println("Failed to get file from Backblaze:", err.Error())
		return
	} else {
		log.Println("Success to get file from Backblaze")
	}

	tmpl, err := template.New("emailTemplate").Parse(htmlFile.String())
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
