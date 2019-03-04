package handler

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/partyzanex/esender/domain"
	"net/http"
)

func (h *Handler) UpdateEmail(ctx echo.Context) error {
	email := domain.Email{}

	if err := ctx.Bind(&email); err != nil {
		return h.errorResponse(ctx, errors.Wrap(err, UnmarshalErr))
	}

	result, err := h.Storage.Update(ctx.Request().Context(), email)
	if err != nil {
		return h.errorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, &Response{
		Data: result,
	})
}
