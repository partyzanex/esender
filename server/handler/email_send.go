package handler

import (
	"github.com/labstack/echo"
	"github.com/partyzanex/esender/domain"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func (h *Handler) EmailSend(ctx echo.Context) error {
	token := ctx.QueryParam("token")

	if !h.Auth(token) {
		return ctx.JSON(http.StatusUnauthorized, &ErrorResponse{
			Error: errors.New("authorization failed"),
		})
	}

	email := &domain.Email{}

	if err := ctx.Bind(email); err != nil {
		return ctx.JSON(http.StatusBadRequest, &ErrorResponse{
			Error: errors.Wrap(err, "unmarshal request failed"),
		})
	}

	email, err := h.Storage.Create(email)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ErrorResponse{
			Error: errors.Wrap(err, "creating email failed"),
		})
	}

	err = h.Sender.Send(email)
	if err != nil {
		email.Status = domain.StatusError
		textError := err.Error()
		email.Error = &textError
	} else {
		email.Status = domain.StatusSent
		now := time.Now()
		email.DTSent = &now
	}

	email, err = h.Storage.Update(email)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ErrorResponse{
			Error: errors.Wrap(err, "updating email failed"),
		})
	}

	return ctx.JSON(http.StatusOK, email)
}
