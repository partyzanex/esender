package domain

import (
	"encoding/base64"
	"net/mail"
	"time"

	"github.com/pkg/errors"
)

const (
	MIME_HTML = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	MIME_TEXT = "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
)

type Email struct {
	ID        int         `json:"id"`
	To        []string    `json:"to"`
	From      string      `json:"from"`
	Subject   string      `json:"subject"`
	Body      string      `json:"body"`
	MimeType  string      `json:"mime_type"`
	Status    EmailStatus `json:"status"`
	Error     *string     `json:"error,omitempty"`
	DTCreated time.Time   `json:"dt_created"`
	DTUpdated *time.Time  `json:"dt_updated"`
	DTSent    *time.Time  `json:"dt_sent"`
}

func (email *Email) GetBody() []byte {
	body, _ := base64.StdEncoding.DecodeString(email.Body)
	return body
}

func (email *Email) GetMimeType() string {
	if email.MimeType == "html" {
		return MIME_HTML
	}

	return MIME_TEXT
}

func (email *Email) Validate() error {
	if len(email.To) == 0 {
		return errors.New("emails recipients list (to) is empty")
	}

	if email.From == "" {
		return errors.New("email sender name (from) is empty")
	}

	if _, err := mail.ParseAddress(email.From); err != nil {
		return errors.Wrap(err, "parsing from address failed")
	}

	for _, to := range email.To {
		if _, err := mail.ParseAddress(to); err != nil {
			return errors.Wrap(err, "parsing to address failed")
		}
	}

	if email.DTCreated.IsZero() {
		return errors.New("dt created must be provided")
	}

	if email.Body == "" {
		return errors.New("email body is empty")
	}

	if !email.Status.IsValid() {
		return errors.New("invalid email status")
	}

	return nil
}

type EmailStatus string

func (status EmailStatus) IsValid() bool {
	switch status {
	case StatusCreated:
		fallthrough
	case StatusSent:
		fallthrough
	case StatusError:
		return true
	}

	return false
}

const (
	StatusCreated EmailStatus = "created"
	StatusSent    EmailStatus = "sent"
	StatusError   EmailStatus = "error"
)

type Sender interface {
	Send(email *Email) error
}

type Filter struct {
	To, From      string
	Status        EmailStatus
	After, Before *time.Time
	Limit, Offset int
}

type EmailStorage interface {
	List(filter *Filter) ([]*Email, error)
	Create(email *Email) (*Email, error)
	Update(email *Email) (*Email, error)
}
