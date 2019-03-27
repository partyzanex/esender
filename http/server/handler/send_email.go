package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/partyzanex/esender/domain"
	"github.com/pkg/errors"
)

func (h *Handler) SendEmail(ctx echo.Context) error {
	var err error

	email := domain.Email{}

	if err := ctx.Bind(&email); err != nil {
		return h.errorResponse(ctx, errors.Wrap(err, UnmarshalErr))
	}

	var sender domain.EmailSender

	if name := ctx.QueryParam("sender"); name != "" {
		var ok bool
		sender, ok = h.Senders.Get(name)
		if !ok {
			return h.errorResponse(ctx, errors.New(UndefinedSenderErr))
		}
	} else {
		sender, _ = h.Senders.Get("")
	}

	result, err := h.Storage.Create(ctx.Request().Context(), email)
	if err != nil {
		return h.errorResponse(ctx, errors.Wrap(err, "creating email failed"))
	}
	defer func() {
		var errS error

		result, errS = h.Storage.Update(ctx.Request().Context(), *result)
		if errS != nil {
			if err != nil {
				err = errors.Wrap(errS, err.Error())
			} else {
				err = errS
			}
		}
	}()

	err = sender.Send(*result)
	if err != nil {
		errStr := err.Error()
		result.Error = &errStr
		result.Status = domain.StatusError

		return ctx.JSON(http.StatusInternalServerError, &Response{
			Error: errors.Wrap(err, domain.SendEmailErr),
			Data:  result,
		})
	}

	DTSent := time.Now()
	result.DTSent = &DTSent
	result.Status = domain.StatusSent

	return ctx.JSON(h.getStatusCode(err), &Response{
		Error: err,
		Data:  result,
	})
}
