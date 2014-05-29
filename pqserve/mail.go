package main

import (
	"fmt"
	"net/smtp"
	"strings"
)

func sendmail(to, subject, body string) (err error) {
	msg := fmt.Sprintf(`From: PaQu <%s>
To: %s
Subject: %s

%s
`, Cfg.Mailfrom, to, subject, body)

	if Cfg.Smtpuser != "" {
		auth := smtp.PlainAuth("", Cfg.Smtpuser, Cfg.Smtppass, strings.Split(Cfg.Smtpserv, ":")[0])
		err = smtp.SendMail(Cfg.Smtpserv, auth, Cfg.Mailfrom, []string{to}, []byte(msg))
	} else {
		err = smtp.SendMail(Cfg.Smtpserv, nil, Cfg.Mailfrom, []string{to}, []byte(msg))
	}
	if err != nil {
		logerr(err)
	}
	return
}
