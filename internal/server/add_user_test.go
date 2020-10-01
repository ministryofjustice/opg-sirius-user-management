package server

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockAddUserClient struct {
	count            int
	lastCookies      []*http.Cookie
	lastEmail        string
	lastFirstname    string
	lastSurname      string
	lastOrganisation string
	lastRoles        []string
	err              error
}

func (m *mockAddUserClient) AddUser(ctx context.Context, cookies []*http.Cookie, email, firstname, surname, organisation string, roles []string) error {
	m.count += 1
	m.lastCookies = cookies
	m.lastEmail = email
	m.lastFirstname = firstname
	m.lastSurname = surname
	m.lastOrganisation = organisation
	m.lastRoles = roles

	return m.err
}

func TestGetAddUser(t *testing.T) {
	assert := assert.New(t)

	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := addUser(nil, nil, template, "http://sirius")(w, r)
	assert.Nil(err)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addUserVars{
		Path:      "/path",
		SiriusURL: "http://sirius",
	}, template.lastVars)
}

func TestPostAddUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddUserClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("email=a&firstname=b&surname=c&organisation=d&roles=e&roles=f"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := addUser(nil, client, template, "http://sirius")(w, r)
	assert.Equal(RedirectError("/users"), err)

	assert.Equal(r.Cookies(), client.lastCookies)
	assert.Equal("a", client.lastEmail)
	assert.Equal("b", client.lastFirstname)
	assert.Equal("c", client.lastSurname)
	assert.Equal("d", client.lastOrganisation)
	assert.Equal([]string{"e", "f"}, client.lastRoles)

	assert.Equal(0, template.count)
}

func TestPostAddUserUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddUserClient{err: sirius.ErrUnauthorized}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := addUser(nil, client, template, "http://sirius")(w, r)
	assert.Equal(sirius.ErrUnauthorized, err)

	assert.Equal(0, template.count)
}

func TestPostAddUserValidationError(t *testing.T) {
	assert := assert.New(t)

	errors := sirius.ValidationErrors{
		"x": {
			"y": "z",
		},
	}
	client := &mockAddUserClient{
		err: sirius.ValidationError{
			Errors: errors,
		},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := addUser(nil, client, template, "http://sirius")(w, r)
	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusBadRequest, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addUserVars{
		Path:      "/path",
		SiriusURL: "http://sirius",
		Errors:    errors,
	}, template.lastVars)
}

func TestPostAddUserOtherError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	logger := log.New(ioutil.Discard, "", 0)
	client := &mockAddUserClient{err: expectedErr}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := addUser(logger, client, template, "http://sirius")(w, r)
	assert.Equal(expectedErr, err)

	assert.Equal(0, template.count)
}
