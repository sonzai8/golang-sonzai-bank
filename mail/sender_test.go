package mail

import (
	"github.com/sonzai8/golang-sonzai-bank/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skip, please run without -short")
	}
	config, err := utils.LoadConfig("..")
	if err != nil {
		t.Errorf("load config failed, err:%v\n", err)
	}
	name := config.EmailConfig.EmailSenderName
	email := config.EmailConfig.EmailSenderAddress
	password := config.EmailConfig.EmailSenderPassword

	sender := NewGmailSender(name, email, password)

	subject := "A Test Email"
	content := `
	<h1> HelloWord </h1>
	<p>This is a test email address from : <a href:="https://sonzai.vn"> Sonzai.vn </a> </p>
	`
	to := []string{"phanvansoninfo@gmail.com"}
	attachFile := []string{"../README.MD"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFile)
	require.NoError(t, err)

}
