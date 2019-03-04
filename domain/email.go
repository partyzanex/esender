package domain

import (
	"encoding/base64"
	"net/mail"
	"time"

	"github.com/pkg/errors"
)

const (
	ParseAddressErr    = "parsing from address failed"
	EmailValidationErr = "email validation failed"
	SearchingEmailErr  = "searching email failed"
	UpdatingEmailErr   = "updating email failed"
	InsertionEmailErr  = "inserting email failed"
	GetEmailBodyErr    = "getting body failed"
	SendEmailErr       = "sending email failed"
)

var (
	ErrEmptyRecipients    = errors.New("emails recipients list (to) is empty")
	ErrEmptySender        = errors.New("email sender name (from) is empty")
	ErrInvalidMimeType    = errors.New("invalid mime type")
	ErrEmptyCreatedDate   = errors.New("dt created must be provided")
	ErrEmptyBody          = errors.New("email body is empty")
	ErrInvalidEmailStatus = errors.New("invalid email status")
	ErrRequiredEmailID    = errors.New("email id is required")
)

type Email struct {
	ID         int           `json:"id"`
	Recipients []string      `json:"recipients"`
	CC         []string      `json:"cc"`
	BCC        []string      `json:"bcc"`
	Sender     string        `json:"from"`
	Subject    string        `json:"subject"`
	Body       string        `json:"body"`
	MimeType   MimeTypeAlias `json:"mime_type"`
	Status     EmailStatus   `json:"status"`
	Error      *string       `json:"error,omitempty"`
	DTCreated  time.Time     `json:"dt_created"`
	DTUpdated  *time.Time    `json:"dt_updated"`
	DTSent     *time.Time    `json:"dt_sent"`
}

func (email *Email) GetBody() ([]byte, error) {
	body, err := base64.StdEncoding.DecodeString(email.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (email *Email) Validate() error {
	if len(email.Recipients) == 0 {
		return ErrEmptyRecipients
	}

	if email.Sender == "" {
		return ErrEmptySender
	}

	if _, err := mail.ParseAddress(email.Sender); err != nil {
		return errors.Wrap(err, ParseAddressErr)
	}

	for _, to := range email.Recipients {
		if _, err := mail.ParseAddress(to); err != nil {
			return errors.Wrap(err, ParseAddressErr)
		}
	}

	if !email.MimeType.IsValid() {
		return ErrInvalidMimeType
	}

	if email.DTCreated.IsZero() {
		return ErrEmptyCreatedDate
	}

	if email.Body == "" {
		return ErrEmptyBody
	}

	if !email.Status.IsValid() {
		return ErrInvalidEmailStatus
	}

	return nil
}
