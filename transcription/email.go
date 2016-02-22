package transcription

import (
	"log"
	"net/smtp"
	"os"
)

func sendEmail(from string, to []string, msg string) {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		os.Getenv("MAIL_EMAIL"),
		os.Getenv("MAIL_PASSWORD"),
		"smtp.gmail.com",
	)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		"smtp.gmail.com:25", auth, from, to, []byte(msg),
	)
	if err != nil {
		log.Fatal(err)
	}
}
