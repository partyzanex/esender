package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func (h *Handler) GetEmail(ctx echo.Context) error {
	var err error
	var emailID int64

	if id := ctx.Param("id"); id != "" {
		emailID, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			return h.errorResponse(ctx, errors.Wrap(err, BadQueryParamErr))
		}
	} else {
		return h.errorResponse(ctx, errors.New(BadQueryParamErr))
	}

	email, err := h.Storage.Get(ctx.Request().Context(), emailID)
	if err != nil {
		return h.errorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, &Response{
		Data: email,
	})
}
