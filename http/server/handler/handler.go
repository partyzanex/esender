package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/partyzanex/esender/domain"
	"github.com/pkg/errors"
)

const (
	UnmarshalErr       = "unmarshal request failed"
	UndefinedSenderErr = "sender not found"
	BadQueryParamErr   = "bad query param"
)

type Handler struct {
	AccessKey string

	Storage domain.EmailStorage
	Senders *domain.Senders
}

func (h *Handler) Auth(key string) bool {
	return key == h.AccessKey
}

type Response struct {
	Error error       `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func (resp Response) MarshalJSON() ([]byte, error) {
	v := struct {
		Error string      `json:"error,omitempty"`
		Data  interface{} `json:"data,omitempty"`
	}{
		Data: resp.Data,
	}

	if resp.Error != nil {
		v.Error = resp.Error.Error()
	}

	b, err := json.Marshal(v)
	if err != nil {
		return nil, errors.Wrap(err, "marshal error")
	}

	return b, nil
}

func (h *Handler) errorResponse(ctx echo.Context, err error) error {
	return ctx.JSON(h.getStatusCode(err), &Response{
		Error: err,
	})
}

func (Handler) getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	code := http.StatusInternalServerError

	switch err {
	case domain.ErrEmptyRecipients:
	case domain.ErrEmptyBody:
	case domain.ErrEmptyCreatedDate:
	case domain.ErrEmptySender:
	case domain.ErrRequiredEmailID:
	case domain.ErrInvalidEmailStatus:
	case domain.ErrInvalidMimeType:
		code = http.StatusBadRequest
	}

	if strings.Contains(err.Error(), domain.ParseAddressErr) ||
		strings.Contains(err.Error(), domain.EmailValidationErr) ||
		strings.Contains(err.Error(), UnmarshalErr) {
		code = http.StatusBadRequest
	}

	if strings.Contains(err.Error(), UndefinedSenderErr) {
		code = http.StatusNotFound
	}

	return code
}
