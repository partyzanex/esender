package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strconv"

	"github.com/partyzanex/esender/domain"
	"github.com/pkg/errors"
)

type Config struct {
	domain.AgentConfig

	Host               string
	Port               uint16
	UserName, Password string
	TLS                bool
}

func (config Config) Address() string {
	return config.Host + ":" + strconv.FormatInt(int64(config.Port), 10)
}

type Sender struct {
	config Config
	auth   smtp.Auth
}

func (Sender) Name() string {
	return "smtp"
}

func (sender *Sender) Send(email domain.Email) (bool, error) {
	if err := email.Validate(); err != nil {
		return false, errors.Wrap(err, domain.EmailValidationErr)
	}

	auth := smtp.PlainAuth("", sender.config.UserName, sender.config.Password, sender.config.Host)

	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: sender.config.TLS,
		ServerName:         sender.config.Host,
	}

	conn, err := tls.Dial("tcp", sender.config.Address(), tlsConfig)
	if err != nil {
		return false, errors.Wrap(err, "dial connection error")
	}

	client, err := smtp.NewClient(conn, sender.config.Host)
	if err != nil {
		return false, errors.Wrap(err, "creating client failed")
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return false, errors.Wrap(err, "auth failed")
	}

	for _, to := range email.Recipients {
		headers := make(map[string]string)
		headers["From"] = email.Sender.String()
		headers["To"] = to.String()
		headers["Subject"] = email.Subject

		message := ""
		for key, value := range headers {
			message += fmt.Sprintf("%s: %s\r\n", key, value)
		}

		message += email.MimeType.Header()
		message += "\r\n" + email.Body

		if err = client.Mail(email.Sender.Address); err != nil {
			return false, errors.Wrap(err, "set sender failed")
		}

		if err = client.Rcpt(to.Address); err != nil {
			return false, errors.Wrap(err, "set recipient failed")
		}

		writer, err := client.Data()
		if err != nil {
			return false, errors.Wrap(err, "getting writer failed")
		}

		_, err = writer.Write([]byte(message))
		if err != nil {
			return false, errors.Wrap(err, "write message failed")
		}

		err = writer.Close()
		if err != nil {
			return false, errors.Wrap(err, "closing writer failed")
		}

		client.Reset()
	}

	return true, nil
}

func (sender *Sender) AgentConfig() domain.AgentConfig {
	return domain.AgentConfig{
		Interval: sender.config.Interval,
		Pause:    sender.config.Interval,
		Limit:    sender.config.Limit,
		Status:   sender.config.Status,
	}
}

func New(config Config) *Sender {
	return &Sender{
		config: config,
		auth:   smtp.PlainAuth("", config.UserName, config.Password, config.Host),
	}
}
