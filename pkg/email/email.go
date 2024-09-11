package email

import (
	"net/smtp"
)

type Sender interface {
	SendMail(to []string, msg []byte) error
}

type SMTP struct {
	auth smtp.Auth

	sender         string
	senderPassword string

	SMTPServer string
	SMTPPort   string
}

func NewSMPTServer(sender string, senderPassword string, smtpServer string, smtpPort string) *SMTP {
	auth := smtp.PlainAuth("", sender, senderPassword, smtpServer)

	return &SMTP{
		auth:           auth,
		sender:         sender,
		senderPassword: senderPassword,
		SMTPServer:     smtpServer,
		SMTPPort:       smtpPort,
	}
}

func (s *SMTP) SendMail(to []string, msg []byte) error {
	err := smtp.SendMail(s.SMTPServer+":"+s.SMTPPort, s.auth, s.sender, to, msg)
	if err != nil {
		return err
	}

	return nil
}
