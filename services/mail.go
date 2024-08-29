package services

import (
	"fmt"
	"net/smtp"
)

type Mail struct {
	auth smtp.Auth
	url  string
}

func NewMail(u string, p string, h string, port string) *Mail {
	return &Mail{
		auth: smtp.PlainAuth("", u, p, h),
		url:  h + ":" + port,
	}
}

func (s *Mail) SendMailSimple(from string, to string, sub string, msg string) error {
	toArray := [1]string{to}
	msgBytes := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"From: %s\r\n"+
			"Cc: test@example.org\r\n"+
			"Subject: %s\r\n\r\n"+
			"%s\r\n",
		to,
		from,
		sub,
		msg,
	))

	return s.SendMailRaw(from, toArray[:], msgBytes)
}

func (s *Mail) SendMailRaw(from string, to []string, msg []byte) error {
	return smtp.SendMail(s.url, s.auth, from, to, msg)
}
