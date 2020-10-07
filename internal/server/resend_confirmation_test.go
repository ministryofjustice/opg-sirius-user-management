package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockResendConfirmationClient struct {
	count       int
	lastCookies []*http.Cookie
	lastEmail   string
	err         error
}

func (m *mockResendConfirmationClient) ResendConfirmation(ctx context.Context, cookies []*http.Cookie, email string) error {
	m.count += 1
	m.lastCookies = cookies
	m.lastEmail = email

	return m.err
}

func TestGetResendConfirmation(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := resendConfirmation(nil, nil, "http://sirius")(w, r)
	assert.Equal(RedirectError("/users"), err)
}

func TestPostResendConfirmation(t *testing.T) {
	assert := assert.New(t)

	client := &mockResendConfirmationClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("email=a&id=b"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := resendConfirmation(client, template, "http://sirius")(w, r)
	assert.Nil(err)

	assert.Equal(1, client.count)
	assert.Equal(r.Cookies(), client.lastCookies)
	assert.Equal("a", client.lastEmail)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(resendConfirmationVars{
		Path:      "/path",
		SiriusURL: "http://sirius",
		ID:        "b",
		Email:     "a",
	}, template.lastVars)
}

func TestPostResendConfirmationError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockResendConfirmationClient{
		err: expectedErr,
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := resendConfirmation(client, nil, "http://sirius")(w, r)
	assert.Equal(expectedErr, err)
}
