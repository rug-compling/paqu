package main

import (
	"fmt"
	"net/smtp"
	"strings"
)

func sendmail(to, subject, body string) (err error) {
	msg := strings.Replace(fmt.Sprintf(`From: PaQu <%s>
To: %s
Subject: %s
Content-type: text/plain; charset=UTF-8

%s
`, Cfg.Mailfrom, to, subject, body), "\n", "\r\n", -1)

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
