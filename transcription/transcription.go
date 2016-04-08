// Package transcription implements functions for the manipulation and
// transcription of audio files.
package transcription

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"os/exec"
	"strings"
)

type transcription struct {
	TextTranscription string
	Metadata          string
}

// SendEmail connects to an email server at host:port, switches to TLS,
// authenticates on TLS connections using the username and password, and sends
// an email from address from, to address to, with subject line subject with
// message body.
func SendEmail(username string, password string, host string, port int, to []string, subject string, body string) error {
	from := username
	auth := smtp.PlainAuth("", username, password, host)

	// The msg parameter should be an RFC 822-style email with headers first,
	// a blank line, and then the message body. The lines of msg should be CRLF
	// terminated.
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

// StartTranscription transcribes a given file using Sphinx.
// File name should be in "name.wav" format.
func StartTranscription(fileName string, command string) error {

	cmd := exec.Command("java", "-jar", command, fileName)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	outputFileName := "/Sphinx/files/" + fileName + "-json.txt"
	//once json is sent somewhere, capture the output
	if _, err := transcriptionOutputToStruct(outputFileName); err != nil {
		return err
	}
	return nil
}

// transcriptionOutputToStruct takes a text file and reads its input
// into a Go struct
func transcriptionOutputToStruct(fileName string) (transcription, error) {
	var jsonData transcription
	file, err := os.Open(fileName)
	r := bufio.NewReader(file)

	bytesText, err := r.ReadBytes('\n')
	if err != nil {
		return jsonData, err
	}
	bytesMeta, err := r.ReadBytes('\n')
	if err != nil || err != io.EOF {
		return jsonData, err
	}

	nText := bytes.IndexByte(bytesText, 0)
	nMeta := bytes.IndexByte(bytesMeta, 0)
	sText := string(bytesText[:nText])
	sMeta := string(bytesText[:nMeta])

	jsonData = transcription{TextTranscription: sText, Metadata: sMeta}
	return jsonData, nil
}

// ConvertAudioIntoWavFormat converts encoded audio into the required format.
func ConvertAudioIntoWavFormat(fn string) error {
	// http://cmusphinx.sourceforge.net/wiki/faq
	// -ar 16000 sets frequency to required 16khz
	// -ac 1 sets the number of audio channels to 1
	cmd := exec.Command("ffmpeg", "-i", fn, "-ar", "16000", "-ac", "1", fn+".wav")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// ConvertAudioIntoFlacFormat converts files into .flac format.
func ConvertAudioIntoFlacFormat(fn string) error {
	// -ar 16000 sets frequency to required 16khz
	// -ac 1 sets the number of audio channels to 1
	cmd := exec.Command("ffmpeg", "-i", fn, "-ar", "16000", "-ac", "1", fn+".flac")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// DownloadFileFromURL locally downloads an audio file stored at url.
func DownloadFileFromURL(url string) error {
	// Taken from https://github.com/thbar/golang-playground/blob/master/download-files.go
	output, err := os.Create(fileNameFromURL(url))
	if err != nil {
		return err
	}
	defer output.Close()

	// Get file contents
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Write the body to file
	_, err = io.Copy(output, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func fileNameFromURL(url string) string {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	return fileName
}

// MakeTaskFunction returns a task function for transcription using transcription functions.
func MakeTaskFunction(audioURL string, emailAddresses []string) func() error {
	return func() error {
		fileName := fileNameFromURL(audioURL)
		if err := DownloadFileFromURL(audioURL); err != nil {
			return err
		}
		if err := ConvertAudioIntoWavFormat(fileName); err != nil {
			return err
		}
		return nil
	}
}
