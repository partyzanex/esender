package handler

import (
	"encoding/json"
	"github.com/partyzanex/esender/domain"
	"github.com/pkg/errors"
)

type Handler struct {
	AuthKey string

	Storage domain.EmailStorage
	Sender  domain.Sender
}

func (h *Handler) Auth(key string) bool {
	return key == h.AuthKey
}

type ErrorResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error,omitempty"`
}

func (e ErrorResponse) MarshalJSON() ([]byte, error) {
	v := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}{
		Success: false,
	}

	if e.Error != nil {
		v.Error = e.Error.Error()
	}

	b, err := json.Marshal(v)
	if err != nil {
		return nil, errors.Wrap(err, "marshal error")
	}

	return b, nil
}
