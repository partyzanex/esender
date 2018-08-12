package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strconv"

	"github.com/partyzanex/esender/domain"
	"github.com/pkg/errors"
)

type SMTPConfig struct {
	Host               string
	Port               uint16
	UserName, Password string
}

func (config SMTPConfig) Address() string {
	return config.Host + ":" + strconv.FormatInt(int64(config.Port), 10)
}

type SMTPSender struct {
	config SMTPConfig
	auth   smtp.Auth
}

func (sender *SMTPSender) Send(email *domain.Email) error {
	if err := email.Validate(); err != nil {
		return errors.Wrap(err, "email validation failed")
	}

	auth := smtp.PlainAuth("", sender.config.UserName, sender.config.Password, sender.config.Host)

	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         sender.config.Host,
	}

	conn, err := tls.Dial("tcp", sender.config.Address(), tlsConfig)
	if err != nil {
		return errors.Wrap(err, "dial connection error")
	}

	client, err := smtp.NewClient(conn, sender.config.Host)
	if err != nil {
		return errors.Wrap(err, "creating client failed")
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return errors.Wrap(err, "auth failed")
	}

	for _, to := range email.To {
		headers := make(map[string]string)
		headers["From"] = email.From
		headers["To"] = to
		headers["Subject"] = email.Subject

		message := ""
		for key, value := range headers {
			message += fmt.Sprintf("%s: %s\r\n", key, value)
		}

		message += email.GetMimeType()
		message += "\r\n" + string(email.GetBody())

		if err = client.Mail(headers["From"]); err != nil {
			return errors.Wrap(err, "set sender failed")
		}

		if err = client.Rcpt(to); err != nil {
			return errors.Wrap(err, "set recipient failed")
		}

		writer, err := client.Data()
		if err != nil {
			return errors.Wrap(err, "getting writer failed")
		}

		_, err = writer.Write([]byte(message))
		if err != nil {
			return errors.Wrap(err, "write message failed")
		}

		err = writer.Close()
		if err != nil {
			return errors.Wrap(err, "closing writer failed")
		}

		client.Reset()
	}

	return nil
}

func NewSMTPSender(config SMTPConfig) *SMTPSender {
	return &SMTPSender{
		config: config,
		auth:   smtp.PlainAuth("", config.UserName, config.Password, config.Host),
	}
}
