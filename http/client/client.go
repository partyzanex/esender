package client

import (
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/resty.v1"
	"github.com/pkg/errors"
	"github.com/partyzanex/esender/domain"
)

const (
	emailPath     = "emails"
	sendEmailPath = "emails/send"

	requestErrTpl = "request returns code: %d with error: %s"
)

type result struct {
	Data  *domain.Email `json:"data,omitempty"`
	Error *string       `json:"error,omitempty"`
}

type Client struct {
	client *resty.Client

	BaseURL string
}

func (c Client) url(path string) string {
	return fmt.Sprintf("%s/%s", c.BaseURL, path)
}

func (c *Client) CreateEmail(email domain.Email) (*domain.Email, error) {
	result := &result{}

	resp, err := c.client.R().
		SetBody(email).
		SetResult(result).
		Post(c.url(emailPath))
	if err != nil {
		return nil, errors.Wrap(err, "creating email failed")
	}

	if code := resp.StatusCode(); code != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf(requestErrTpl, code, *result.Error),
		)
	}

	return result.Data, nil
}

func (c *Client) UpdateEmail(email domain.Email) (*domain.Email, error) {
	result := &result{}

	resp, err := c.client.R().
		SetBody(email).
		SetResult(result).
		Put(c.url(emailPath))
	if err != nil {
		return nil, errors.Wrap(err, "updating email failed")
	}

	if code := resp.StatusCode(); code != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf(requestErrTpl, code, *result.Error),
		)
	}

	return result.Data, nil
}

func (c *Client) SendEmail(email domain.Email) (*domain.Email, error) {
	result := &result{}

	resp, err := c.client.R().
		SetBody(email).
		SetResult(result).
		Post(c.url(sendEmailPath))
	if err != nil {
		return nil, errors.Wrap(err, "sending email failed")
	}

	if code := resp.StatusCode(); code != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf(requestErrTpl, code, *result.Error),
		)
	}

	return result.Data, nil
}

func (c *Client) GetEmail(id int64) (*domain.Email, error) {
	result := &result{}

	resp, err := c.client.R().
		SetResult(result).
		Get(c.url(emailPath) + fmt.Sprintf("/%d", id))
	if err != nil {
		return nil, errors.Wrap(err, "getting email failed")
	}

	if code := resp.StatusCode(); code != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf(requestErrTpl, code, *result.Error),
		)
	}

	return result.Data, nil
}

func (c *Client) SearchEmails(filter *domain.Filter) ([]*domain.Email, error) {
	result := &struct {
		Data  []*domain.Email `json:"data,omitempty"`
		Error *string         `json:"error,omitempty"`
	}{}

	request := c.client.R().
		SetResult(result)

	if filter != nil {
		if filter.Recipient != "" {
			request.SetQueryParam("recipient", filter.Recipient)
		}

		if filter.Sender != "" {
			request.SetQueryParam("sender", filter.Sender)
		}
		if filter.Status != "" {
			request.SetQueryParam("status", filter.Status.String())
		}
		if filter.TimeRange.IsValid() {
			request.SetQueryParam("since", filter.TimeRange.Since().Format(domain.DateTimeLayout))
			request.SetQueryParam("till", filter.TimeRange.Till().Format(domain.DateTimeLayout))
		}
		if filter.Limit > 0 {
			request.SetQueryParam("limit", strconv.FormatInt(int64(filter.Limit), 10))

			if filter.Offset >= 0 {
				request.SetQueryParam("offset", strconv.FormatInt(int64(filter.Offset), 10))
			}
		}
	}

	resp, err := request.Get(c.url(emailPath))
	if err != nil {
		return nil, errors.Wrap(err, "search for emails failed")
	}

	if code := resp.StatusCode(); code != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf(requestErrTpl, code, *result.Error),
		)
	}

	return result.Data, nil
}

func New(baseURL, user, pass string) *Client {
	return &Client{
		BaseURL: baseURL,
		client:  resty.New().SetBasicAuth(user, pass),
	}
}
