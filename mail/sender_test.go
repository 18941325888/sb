package mail

import (
	"testing"

	"github.com/18941325888/sb/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewMailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
				<h1>Hello world</h1>
				<p>This is a test message from s</p>
				`
	to := []string{"290684858@qq.com"}
	attachFiles := []string{"../sqlc.yaml"}
	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
