package domain

import (
	"encoding/base64"
	"net/mail"
	"time"

	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"log"
)

const (
	ParseAddressErr         = "parsing from address failed"
	EmailValidationErr      = "email validation failed"
	SearchingEmailErr       = "searching email failed"
	UpdatingEmailErr        = "updating email failed"
	InsertionEmailErr       = "inserting email failed"
	SendEmailErr            = "sending email failed"
	ValidationRecipientsErr = "validation of recipients failed"
	ValidationBCCErr        = "validation of BCC failed"
	ValidationCCErr         = "validation of CC failed"
)

var (
	ErrEmptyRecipients     = errors.New("emails recipients list (to) is empty")
	ErrEmptySender         = errors.New("email sender name (from) is empty")
	ErrInvalidMimeType     = errors.New("invalid mime type")
	ErrEmptyCreatedDate    = errors.New("dt created must be provided")
	ErrEmptyBody           = errors.New("email body is empty")
	ErrInvalidEmailStatus  = errors.New("invalid email status")
	ErrRequiredEmailID     = errors.New("email id is required")
	ErrInvalidEmailAddress = errors.New("invalid email address")
)

type Address struct {
	Name    string
	Address string
}

func (address Address) String() string {
	tpl := "%s"
	args := []interface{}{address.Address}

	if address.Name != "" {
		tpl = "\"%s\" <%s>"
		args = []interface{}{address.Name, address.Address}
	}

	return fmt.Sprintf(tpl, args...)
}

func (address Address) Validate() error {
	if !govalidator.IsEmail(address.Address) {
		return ErrInvalidEmailAddress
	}

	return nil
}

func ParseAddress(str string) (*Address, error) {
	address, err := mail.ParseAddress(str)
	if err != nil {
		return nil, err
	}

	return &Address{
		Name:    address.Name,
		Address: address.Address,
	}, nil
}

type Email struct {
	ID         int64         `json:"id"`
	Recipients []Address     `json:"recipients"`
	CC         []Address     `json:"cc"`
	BCC        []Address     `json:"bcc"`
	Sender     Address       `json:"from"`
	Subject    string        `json:"subject"`
	Body       string        `json:"body"`
	MimeType   MimeTypeAlias `json:"mime_type"`
	Status     EmailStatus   `json:"status"`
	Error      *string       `json:"error,omitempty"`
	DTCreated  time.Time     `json:"dt_created"`
	DTUpdated  *time.Time    `json:"dt_updated,omitempty"`
	DTSent     *time.Time    `json:"dt_sent,omitempty"`
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

	if err := email.Sender.Validate(); err != nil {
		return err
	}

	for _, to := range email.Recipients {
		if err := to.Validate(); err != nil {
			return errors.Wrap(err, ValidationRecipientsErr)
		}
	}

	if len(email.BCC) > 0 {
		for _, bcc := range email.BCC {
			log.Println("bcc", bcc)
			if err := bcc.Validate(); err != nil {
				return errors.Wrap(err, ValidationBCCErr)
			}
		}
	}

	if len(email.CC) > 0 {
		for _, cc := range email.CC {
			if err := cc.Validate(); err != nil {
				return errors.Wrap(err, ValidationCCErr)
			}
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
