package services

import (
	"database/sql"
	"strings"
	"time"

	"github.com/partyzanex/esender/domain"
	"github.com/partyzanex/esender/models/mysql"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

type EmailMySqlStorage struct {
	DB *sql.DB
}

func (storage *EmailMySqlStorage) List(filter *domain.Filter) ([]*domain.Email, error) {
	var mods []qm.QueryMod

	if filter != nil {
		if filter.To != "" {
			mods = append(mods, qm.Where("to like ?", "%"+filter.To+"%"))
		}

		if filter.From != "" {
			mods = append(mods, qm.Where("from like ?", "%"+filter.From+"%"))
		}

		if filter.Status != "" {
			if !filter.Status.IsValid() {
				return nil, errors.New("invalid status")
			}

			mods = append(mods, qm.Where("status = ?", filter.Status))
		}

		if filter.Limit > 0 {
			mods = append(mods, qm.Limit(filter.Limit))
		}

		if filter.Offset >= 0 {
			mods = append(mods, qm.Offset(filter.Offset))
		}

		if filter.After != nil && !filter.After.IsZero() {
			mods = append(mods, qm.Where("created_at >= ?", *filter.After))
		}

		if filter.Before != nil && !filter.Before.IsZero() {
			mods = append(mods, qm.Where("created_at <= ?", *filter.Before))
		}
	}

	models, err := mysql.Emails(storage.DB, mods...).All()
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "searching emails failed")
	}

	emails := make([]*domain.Email, len(models))

	for i, model := range models {
		emails[i] = emailFromModel(model)
	}

	return emails, nil
}

func (storage *EmailMySqlStorage) Create(email *domain.Email) (*domain.Email, error) {
	if email == nil {
		return nil, errors.New("email is empty")
	}

	model := &mysql.Email{
		To:        strings.Join(email.To, ";"),
		From:      email.From,
		Text:      email.Body,
		Title:     email.Subject,
		MimeType:  email.MimeType,
		Status:    string(email.Status),
		DTCreated: int(time.Now().Unix()),
	}

	err := model.Insert(storage.DB)
	if err != nil {
		return nil, errors.Wrap(err, "inserting email failed")
	}

	return emailFromModel(model), nil
}

func (storage *EmailMySqlStorage) Update(email *domain.Email) (*domain.Email, error) {
	if email == nil {
		return nil, errors.New("email is empty")
	}

	if email.ID == 0 {
		return nil, errors.New("email id is required")
	}

	model, err := mysql.Emails(storage.DB, qm.Where("id = ?", email.ID)).One()
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "searching email failed")
	}

	model.Status = string(email.Status)
	model.To = strings.Join(email.To, ";")
	model.From = email.From
	model.Title = email.Subject
	model.MimeType = email.MimeType

	if email.Error != nil {
		model.Error.Valid = true
		model.Error.String = *email.Error
	}

	model.DTUpdated.Valid = true
	model.DTUpdated.Int = int(time.Now().Unix())

	if email.DTSent != nil {
		model.DTSent.Valid = true
		model.DTSent.Int = int(email.DTSent.Unix())
	}

	err = model.Update(storage.DB)
	if err != nil {
		return nil, errors.Wrap(err, "updating email failed")
	}

	return emailFromModel(model), nil
}

func emailFromModel(model *mysql.Email) *domain.Email {
	email := &domain.Email{
		ID:        model.ID,
		To:        strings.Split(model.To, ";"),
		From:      model.From,
		Subject:   model.Title,
		MimeType:  model.MimeType,
		Body:      model.Text,
		Status:    domain.EmailStatus(model.Status),
		DTCreated: time.Unix(int64(model.DTCreated), 0),
	}

	if model.Error.Valid {
		err := model.Error.String
		email.Error = &err
	}

	if model.DTUpdated.Valid {
		DTUpdated := time.Unix(int64(model.DTUpdated.Int), 0)
		email.DTUpdated = &DTUpdated
	}

	if model.DTSent.Valid {
		DTSent := time.Unix(int64(model.DTSent.Int), 0)
		email.DTSent = &DTSent
	}

	return email
}
