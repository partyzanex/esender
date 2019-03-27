package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/partyzanex/esender/domain"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"strconv"
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

func (res result) error() string {
	respErr := ""
	if res.Error != nil {
		respErr = *res.Error
	}

	return respErr
}

type Client struct {
	client *resty.Client

	BaseURL string
}

func (c Client) url(path string) string {
	return fmt.Sprintf("%s/%s", c.BaseURL, path)
}

func (c *Client) CreateEmail(email domain.Email) (*domain.Email, error) {
	resp, err := c.client.R().
		SetBody(email).
		Post(c.url(emailPath))
	if err != nil {
		return nil, errors.Wrap(err, "creating email failed")
	}

	result, err := c.bindResult(resp.Body())
	if err != nil {
		return nil, err
	}

	if code := resp.StatusCode(); code != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf(requestErrTpl, code, result.error()),
		)
	}

	return result.Data, nil
}

func (Client) bindResult(body []byte) (*result, error) {
	result := &result{}

	dec := json.NewDecoder(bytes.NewReader(body))

	err := dec.Decode(result)
	if err != nil {
		return nil, errors.Wrap(err, "decoding json failed")
	}

	return result, nil
}

func (c *Client) UpdateEmail(email domain.Email) (*domain.Email, error) {
	resp, err := c.client.R().
		SetBody(email).
		Put(c.url(emailPath))
	if err != nil {
		return nil, errors.Wrap(err, "updating email failed")
	}

	result, err := c.bindResult(resp.Body())
	if err != nil {
		return nil, err
	}

	if code := resp.StatusCode(); code != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf(requestErrTpl, code, result.error()),
		)
	}

	return result.Data, nil
}

func (c *Client) SendEmail(email domain.Email) (*domain.Email, error) {
	resp, err := c.client.R().
		SetBody(email).
		Post(c.url(sendEmailPath))
	if err != nil {
		return nil, errors.Wrap(err, "sending email failed")
	}

	result, err := c.bindResult(resp.Body())
	if err != nil {
		return nil, err
	}

	if code := resp.StatusCode(); code != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf(requestErrTpl, code, result.error()),
		)
	}

	return result.Data, nil
}

func (c *Client) GetEmail(id int64) (*domain.Email, error) {
	resp, err := c.client.R().
		Get(c.url(emailPath) + fmt.Sprintf("/%d", id))
	if err != nil {
		return nil, errors.Wrap(err, "getting email failed")
	}

	result, err := c.bindResult(resp.Body())
	if err != nil {
		return nil, err
	}

	if code := resp.StatusCode(); code != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf(requestErrTpl, code, result.error()),
		)
	}

	return result.Data, nil
}

func (c *Client) SearchEmails(filter *domain.Filter) ([]*domain.Email, error) {
	request := c.client.R()

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

	result := &struct {
		Data  []*domain.Email `json:"data,omitempty"`
		Error *string         `json:"error,omitempty"`
	}{}

	dec := json.NewDecoder(bytes.NewReader(resp.Body()))

	err = dec.Decode(result)
	if err != nil {
		return nil, errors.Wrap(err, "decoding json failed")
	}

	if code := resp.StatusCode(); code != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf(requestErrTpl, code, result.Error),
		)
	}

	return result.Data, nil
}

func New(baseURL, user, pass string) *Client {
	client := resty.New().
		SetBasicAuth(user, pass).
		SetHeader("Content-Type", "application/json")

	return &Client{
		BaseURL: baseURL,
		client:  client,
	}
}
