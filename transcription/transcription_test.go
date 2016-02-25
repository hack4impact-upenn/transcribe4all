package transcription

import (
	"net/smtp" // mock
	"testing"

	"github.com/golang/mock/gomock"
)

func TestSendEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup the mock package
	smtp.MOCK().SetController(ctrl)

	username := "test@email.com"
	password := "123456"
	host := "smtp.gmail.com"
	port := 25
	addr := host + ":" + string(port)
	to := []string{"to@email.com"}
	subject := "subject"
	body := "body"

	gomock.InOrder(
		smtp.EXPECT().PlainAuth("", username, password, "smtp.gmail.com").Times(1),
		smtp.EXPECT().SendMail(addr, gomock.Any(), username, to, gomock.Any()).Times(1),
	)

	SendEmail(username, password, host, port, to, subject, body)
}
