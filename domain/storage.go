package domain

import (
	"time"
	"context"
)

const (
	DateCreated TimeRangeProp = "created"
	DateUpdated TimeRangeProp = "updated"
	DateSent    TimeRangeProp = "sent"

	DateTimeLayout = "2006-01-02 15:45:06"
)

type TimeRangeProp string

func (trp TimeRangeProp) IsValid() bool {
	return trp == DateCreated || trp == DateUpdated || trp == DateSent
}

type TimeRange struct {
	prop  TimeRangeProp
	dates [2]time.Time
}

func (tr *TimeRange) SetProp(prop TimeRangeProp) {
	tr.prop = prop
}

func (tr *TimeRange) Prop() TimeRangeProp {
	return tr.prop
}

func (tr *TimeRange) SetSince(dt time.Time) {
	tr.dates[0] = dt
}

func (tr *TimeRange) SetTill(dt time.Time) {
	tr.dates[1] = dt
}

func (tr *TimeRange) Since() time.Time {
	return tr.dates[0]
}

func (tr *TimeRange) Till() time.Time {
	return tr.dates[1]
}

func (tr *TimeRange) IsValid() bool {
	return !tr.dates[0].IsZero() && !tr.dates[1].IsZero() && tr.prop.IsValid()
}

type Filter struct {
	Recipient, Sender string
	Status            EmailStatus
	TimeRange         TimeRange
	Limit, Offset     int
}

type EmailStorage interface {
	Search(ctx context.Context, filter *Filter) ([]*Email, error)
	Get(ctx context.Context, id int64) (*Email, error)
	Create(ctx context.Context, email Email) (*Email, error)
	Update(ctx context.Context, email Email) (*Email, error)
}
