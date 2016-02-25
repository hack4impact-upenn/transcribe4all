package transcription

import (
	"log"
	"net/smtp"
	"os"
	"strings"
)

// SendEmail authenticates from environment variables and sends email using arguments.
func SendEmail(from string, to []string, subject string, body string) {
	// Set up authentication information.
	email := os.Getenv("MAIL_EMAIL")
	password := os.Getenv("MAIL_PASSWORD")
	auth := smtp.PlainAuth("", email, password, "smtp.gmail.com")

	msg := []byte("To: " + strings.Join(to, ", ") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")
	err := smtp.SendMail("smtp.gmail.com:25", auth, from, to, msg)

	if err != nil {
		log.Fatal(err)
	}
}
