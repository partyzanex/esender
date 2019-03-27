package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/partyzanex/esender/domain"
	"github.com/pkg/errors"
)

func (h *Handler) SearchEmails(ctx echo.Context) error {
	filter := &domain.Filter{
		Status:    domain.EmailStatus(ctx.QueryParam("status")),
		Sender:    ctx.QueryParam("sender"),
		Recipient: ctx.QueryParam("recipient"),
	}

	if lim := ctx.QueryParam("limit"); lim != "" {
		limit, err := strconv.ParseInt(lim, 10, 64)
		if err != nil {
			return h.errorResponse(ctx, errors.Wrap(err, BadQueryParamErr))
		}

		filter.Limit = int(limit)
	}

	if offs := ctx.QueryParam("offset"); offs != "" {
		offset, err := strconv.ParseInt(offs, 10, 64)
		if err != nil {
			return h.errorResponse(ctx, errors.Wrap(err, BadQueryParamErr))
		}

		filter.Offset = int(offset)
	}

	if prop := ctx.QueryParam("dt"); prop != "" {
		filter.TimeRange.SetProp(domain.TimeRangeProp(prop))
	}

	if dt := ctx.QueryParam("since"); dt != "" {
		since, err := time.Parse(domain.DateTimeLayout, dt)
		if err != nil {
			return h.errorResponse(ctx, errors.Wrap(err, BadQueryParamErr))
		}

		filter.TimeRange.SetSince(since)
	}

	if dt := ctx.QueryParam("till"); dt != "" {
		till, err := time.Parse(domain.DateTimeLayout, dt)
		if err != nil {
			return h.errorResponse(ctx, errors.Wrap(err, BadQueryParamErr))
		}

		filter.TimeRange.SetTill(till)
	}

	emails, err := h.Storage.Search(ctx.Request().Context(), filter)
	if err != nil {
		return h.errorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, &Response{
		Data: emails,
	})
}
