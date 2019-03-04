package mysql

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/partyzanex/esender/boiler/models/mysql"
	"github.com/partyzanex/esender/domain"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

const (
	AddressSeparator = ";"
)

type emailStorage struct {
	db *sql.DB
}

func (storage *emailStorage) Search(ctx context.Context, filter *domain.Filter) ([]*domain.Email, error) {
	var mods []qm.QueryMod

	if filter != nil {
		if filter.Recipient != "" {
			expr := "%"+filter.Recipient+"%"
			clause := "recipients like ? or cc like ? or bcc like ?"
			mods = append(mods, qm.Where(clause, expr, expr, expr))
		}

		if filter.Sender != "" {
			mods = append(mods, qm.Where("sender like ?", "%"+filter.Sender+"%"))
		}

		if filter.Status != "" {
			if !filter.Status.IsValid() {
				return nil, domain.ErrInvalidEmailStatus
			}

			mods = append(mods, qm.Where("status = ?", filter.Status))
		}

		if filter.Limit > 0 {
			mods = append(mods, qm.Limit(filter.Limit))
		}

		if filter.Offset >= 0 {
			mods = append(mods, qm.Offset(filter.Offset))
		}

		if filter.TimeRange.IsValid() {
			dateField := ""
			switch filter.TimeRange.Prop() {
			case domain.DateCreated:
				dateField = "dt_created"
			case domain.DateUpdated:
				dateField = "dt_updated"
			case domain.DateSent:
				dateField = "dt_sent"
			}

			clause := fmt.Sprintf("%s between ? and ?", dateField)
			mods = append(mods, qm.Where(clause, filter.TimeRange.Since(), filter.TimeRange.Till()))
		}
	}

	models, err := mysql.Emails(mods...).All(ctx, storage.db)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, domain.SearchingEmailErr)
	}

	emails := make([]*domain.Email, len(models))

	for i, model := range models {
		emails[i] = emailFromModel(model)
	}

	return emails, nil
}

func (storage *emailStorage) Get(ctx context.Context, id int64) (*domain.Email, error) {
	model, err := mysql.Emails(
		qm.Where("id = ?", id),
	).One(ctx, storage.db)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, domain.SearchingEmailErr)
	}

	return emailFromModel(model), nil
}

func (storage *emailStorage) Create(ctx context.Context, email domain.Email) (*domain.Email, error) {
	if err := email.Validate(); err != nil {
		return nil, errors.Wrap(err, domain.EmailValidationErr)
	}

	model := &mysql.Email{
		Recipients: strings.Join(email.Recipients, AddressSeparator),
		CC:         strings.Join(email.CC, AddressSeparator),
		BCC:        strings.Join(email.BCC, AddressSeparator),
		Sender:     email.Sender,
		Body:       base64.StdEncoding.EncodeToString([]byte(email.Body)),
		Subject:    email.Subject,
		MimeType:   email.MimeType.String(),
		Status:     email.Status.String(),
		DTCreated:  email.DTCreated,
	}

	err := model.Insert(ctx, storage.db, boil.Infer())
	if err != nil {
		return nil, errors.Wrap(err, domain.InsertionEmailErr)
	}

	return emailFromModel(model), nil
}

func (storage *emailStorage) Update(ctx context.Context, email domain.Email) (*domain.Email, error) {
	if email.ID == 0 {
		return nil, domain.ErrRequiredEmailID
	}

	model, err := mysql.Emails(qm.Where("id = ?", email.ID)).One(ctx, storage.db)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, domain.SearchingEmailErr)
	}

	model.Status = string(email.Status)
	model.Recipients = strings.Join(email.Recipients, ";")
	model.CC = strings.Join(email.Recipients, ";")
	model.BCC = strings.Join(email.BCC, ";")
	model.Sender = email.Sender
	model.Subject = email.Subject
	model.MimeType = email.MimeType.String()
	model.Body = base64.StdEncoding.EncodeToString([]byte(email.Body))

	if email.Error != nil {
		model.Error.Valid = true
		model.Error.String = *email.Error
	}

	model.DTUpdated.Valid = true
	model.DTUpdated.Time = time.Now()

	if email.DTSent != nil {
		model.DTSent.Valid = true
		model.DTSent.Time = *email.DTSent
	}

	_, err = model.Update(ctx, storage.db, boil.Infer())
	if err != nil {
		return nil, errors.Wrap(err, domain.UpdatingEmailErr)
	}

	return emailFromModel(model), nil
}

func emailFromModel(model *mysql.Email) *domain.Email {
	body, _ := base64.StdEncoding.DecodeString(model.Body)

	email := &domain.Email{
		ID:         model.ID,
		Recipients: strings.Split(model.Recipients, AddressSeparator),
		CC:         strings.Split(model.CC, AddressSeparator),
		BCC:        strings.Split(model.BCC, AddressSeparator),
		Sender:     model.Sender,
		Subject:    model.Subject,
		MimeType:   domain.MimeTypeAlias(model.MimeType),
		Body:       string(body),
		Status:     domain.EmailStatus(model.Status),
		DTCreated:  model.DTCreated,
	}

	if model.Error.Valid {
		err := model.Error.String
		email.Error = &err
	}

	if model.DTUpdated.Valid {
		email.DTUpdated = &model.DTUpdated.Time
	}

	if model.DTSent.Valid {
		email.DTSent = &model.DTSent.Time
	}

	return email
}

var emailOnce sync.Once

var emailInstance *emailStorage

func EmailStorage(db *sql.DB) *emailStorage {
	emailOnce.Do(func() {
		emailInstance = &emailStorage{
			db: db,
		}
	})

	return emailInstance
}
