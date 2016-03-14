package transcription

import (
	"io"
	"net/smtp"
	"os"
	"os/exec"
	"strings"
)

// SendEmail connects to an email server at host:port, switches to TLS,
// authenticates on TLS connections using the username and password, and sends
// an email from address from, to address to, with subject line subject with message body.
func SendEmail(username string, password string, host string, port int, to []string, subject string, body string) error {
	from := username
	auth := smtp.PlainAuth("", username, password, host)

	// The msg parameter should be an RFC 822-style email with headers first,
	// a blank line, and then the message body. The lines of msg should be CRLF terminated.
	msg := []byte(msgHeaders(from, to, subject) + "\r\n" + body + "\r\n")
	addr := host + ":" + string(port)
	if err := smtp.SendMail(addr, auth, from, to, msg); err != nil {
		return err
	}
	return nil
}

func msgHeaders(from string, to []string, subject string) string {
	fromHeader := "From: " + from
	toHeader := "To: " + strings.Join(to, ", ")
	subjectHeader := "Subject: " + subject
	msgHeaders := []string{fromHeader, toHeader, subjectHeader}
	return strings.Join(msgHeaders, "\r\n")
}

// Transcribe a given file using a Sphinx jar. The fileName should be in
// "name.wav" format and the jarName should be in "name.jar" format.
func startTranscription(fileName string, jarName string) error {

	cmd := exec.Command("java", "-jar", jarName, fileName)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	//Currently copying output to Stdout
	if _, err := io.Copy(os.Stdout, stdout); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
