package transcription

import (
	"net/smtp"
	"strings"
)

// SendEmail authenticates, creates message, and sends email from arguments.
func SendEmail(username string, password string, host string, port int, to []string, subject string, body string) error {
	auth := smtp.PlainAuth("", username, password, host)

	// The msg headers should usually include fields such as "From", "To", "Subject", and "Cc".
	// Sending "Bcc" messages is accomplished by including an email address in
	// the to parameter but not including it in the msg headers.
	from := username
	fromHeader := "From: " + from
	toHeader := "To: " + strings.Join(to, ", ")
	subjectHeader := "Subject: " + subject
	msgHeaders := []string{fromHeader, toHeader, subjectHeader}

	// The msg parameter should be an RFC 822-style email with headers first,
	// a blank line, and then the message body. The lines of msg should be CRLF terminated.
	msg := []byte(strings.Join(msgHeaders, "\r\n") +
		"\r\n" +
		body + "\r\n")
	addr := host + ":" + string(port)
	if err := smtp.SendMail(addr, auth, from, to, msg); err != nil {
		return err
	}
	return nil
}
