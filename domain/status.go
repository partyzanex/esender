package domain

const (
	StatusCreated EmailStatus = "created"
	StatusSent    EmailStatus = "sent"
	StatusError   EmailStatus = "error"
)

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

func (status EmailStatus) String() string {
	return string(status)
}