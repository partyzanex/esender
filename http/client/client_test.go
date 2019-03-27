package client_test

import (
	"os"
	"testing"
	"time"

	"github.com/partyzanex/esender/domain"
	"github.com/partyzanex/esender/http/client"

	test "github.com/google/gxui/testing"
)

func TestClient_CreateEmail(t *testing.T) {
	c := getTestClient()

	email := getTestEmail()

	result, err := c.CreateEmail(email)
	if err != nil {
		t.Fatalf("creating email failed: %s", err)
	}
	if result == nil {
		t.Fatalf("result is nil")
	}

	test.AssertEqual(t, "Sender", result.Sender.String(), email.Sender.String())
	test.AssertEqual(t, "Recipients", result.Recipients, email.Recipients)
	test.AssertEqual(t, "CC", result.CC, email.CC)
	test.AssertEqual(t, "BCC", result.BCC, email.BCC)
	test.AssertEqual(t, "Body", result.Body, email.Body)
	test.AssertEqual(t, "Status", result.Status, email.Status)
	test.AssertEqual(t, "Subject", result.Subject, email.Subject)
	test.AssertEqual(t, "DTCreated", result.DTCreated.Unix(), email.DTCreated.Unix())
	test.AssertEqual(t, "MimeType", result.MimeType, email.MimeType)

	email = *result
	email.Subject = "Test Subject 2"

	result, err = c.UpdateEmail(email)
	if err != nil {
		t.Fatalf("updationg email failed: %s", err)
	}
	if result == nil {
		t.Fatalf("result is nil")
	}

	test.AssertEqual(t, "Sender", result.Sender, email.Sender)
	test.AssertEqual(t, "Recipients", result.Recipients, email.Recipients)
	test.AssertEqual(t, "CC", result.CC, email.CC)
	test.AssertEqual(t, "BCC", result.BCC, email.BCC)
	test.AssertEqual(t, "Body", result.Body, email.Body)
	test.AssertEqual(t, "Status", result.Status, email.Status)
	test.AssertEqual(t, "Subject", result.Subject, email.Subject)
	test.AssertEqual(t, "DTCreated", result.DTCreated.Unix(), email.DTCreated.Unix())
	test.AssertEqual(t, "MimeType", result.MimeType, email.MimeType)

	result2, err := c.GetEmail(int64(result.ID))
	if err != nil {
		t.Fatalf("getting email failed: %s", err)
	}
	if result2 == nil {
		t.Fatalf("result is nil")
	}

	test.AssertEqual(t, "Sender", result.Sender, result2.Sender)
	test.AssertEqual(t, "Recipients", result.Recipients, result2.Recipients)
	test.AssertEqual(t, "CC", result.CC, result2.CC)
	test.AssertEqual(t, "BCC", result.BCC, result2.BCC)
	test.AssertEqual(t, "Body", result.Body, result2.Body)
	test.AssertEqual(t, "Status", result.Status, result2.Status)
	test.AssertEqual(t, "Subject", result.Subject, result2.Subject)
	test.AssertEqual(t, "DTCreated", result.DTCreated.Unix(), result2.DTCreated.Unix())
	test.AssertEqual(t, "MimeType", result.MimeType, result2.MimeType)
}

func TestClient_SendEmail(t *testing.T) {
	c := getTestClient()

	email := getTestEmail()

	result, err := c.SendEmail(email)
	if err != nil {
		t.Fatalf("sending email failed: %s", err)
	}
	if result == nil {
		t.Fatalf("result is nil")
	}

	test.AssertEqual(t, "Sender", result.Sender.String(), email.Sender.String())
	test.AssertEqual(t, "Recipients", result.Recipients, email.Recipients)
	test.AssertEqual(t, "CC", result.CC, email.CC)
	test.AssertEqual(t, "BCC", result.BCC, email.BCC)
	test.AssertEqual(t, "Body", result.Body, email.Body)
	test.AssertEqual(t, "Status", result.Status, domain.StatusSent)
}

func getTestClient() *client.Client {
	return client.New(
		os.Getenv("TEST_BASE_URL"),
		os.Getenv("TEST_USER"),
		os.Getenv("TEST_PASS"),
	)
}

func getTestEmail() domain.Email {
	recipient1 := domain.Address{
		Name:    "John Doe1",
		Address: "johndoe1@example.com",
	}
	recipient2 := domain.Address{
		Name:    "John Doe2",
		Address: "johndoe2@example.com",
	}
	sender := domain.Address{
		Name:    "John Doe",
		Address: "johndoe@example.com",
	}

	return domain.Email{
		Recipients: []domain.Address{
			recipient1, recipient2,
		},
		Sender:    sender,
		MimeType:  domain.HtmlMimeType,
		Subject:   "Test subject",
		Body:      `<h3>Test Heading</h3><p>test message text</p>`,
		Status:    domain.StatusCreated,
		DTCreated: time.Now(),
	}
}
