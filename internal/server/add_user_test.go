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
	addUser struct {
		count            int
		lastCookies      []*http.Cookie
		lastEmail        string
		lastFirstname    string
		lastSurname      string
		lastOrganisation string
		lastRoles        []string
		err              error
	}

	myDetails struct {
		count       int
		lastCookies []*http.Cookie
		err         error
		roles       []string
	}
}

func (m *mockAddUserClient) AddUser(ctx context.Context, cookies []*http.Cookie, email, firstname, surname, organisation string, roles []string) error {
	m.addUser.count += 1
	m.addUser.lastCookies = cookies
	m.addUser.lastEmail = email
	m.addUser.lastFirstname = firstname
	m.addUser.lastSurname = surname
	m.addUser.lastOrganisation = organisation
	m.addUser.lastRoles = roles

	return m.addUser.err
}

func (m *mockAddUserClient) MyDetails(ctx context.Context, cookies []*http.Cookie) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCookies = cookies

	return sirius.MyDetails{Roles: m.myDetails.roles}, m.myDetails.err
}

func TestGetAddUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddUserClient{}
	client.myDetails.roles = []string{"System Admin"}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := addUser(nil, client, template, "http://sirius")(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(0, client.addUser.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addUserVars{
		Path:      "/path",
		SiriusURL: "http://sirius",
	}, template.lastVars)
}

func TestGetAddUserMissingRole(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddUserClient{}
	client.myDetails.roles = []string{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := addUser(nil, client, template, "http://sirius")(w, r)

	assert.Equal(StatusError(http.StatusForbidden), err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(0, client.addUser.count)
	assert.Equal(0, template.count)
}

func TestPostAddUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddUserClient{}
	client.myDetails.roles = []string{"System Admin"}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("email=a&firstname=b&surname=c&organisation=d&roles=e&roles=f"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := addUser(nil, client, template, "http://sirius")(w, r)
	assert.Equal(RedirectError("/users"), err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(r.Cookies(), client.myDetails.lastCookies)

	assert.Equal(1, client.addUser.count)
	assert.Equal(r.Cookies(), client.addUser.lastCookies)
	assert.Equal("a", client.addUser.lastEmail)
	assert.Equal("b", client.addUser.lastFirstname)
	assert.Equal("c", client.addUser.lastSurname)
	assert.Equal("d", client.addUser.lastOrganisation)
	assert.Equal([]string{"e", "f"}, client.addUser.lastRoles)

	assert.Equal(0, template.count)
}

func TestPostAddUserMissingRole(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddUserClient{}
	client.myDetails.roles = []string{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("email=a&firstname=b&surname=c&organisation=d&roles=e&roles=f"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := addUser(nil, client, template, "http://sirius")(w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(r.Cookies(), client.myDetails.lastCookies)

	assert.Equal(0, client.addUser.count)
	assert.Equal(0, template.count)
}

func TestPostAddUserValidationError(t *testing.T) {
	assert := assert.New(t)

	errors := sirius.ValidationErrors{
		"x": {
			"y": "z",
		},
	}
	client := &mockAddUserClient{}
	client.myDetails.roles = []string{"System Admin"}
	client.addUser.err = sirius.ValidationError{
		Errors: errors,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := addUser(nil, client, template, "http://sirius")(w, r)
	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusBadRequest, resp.StatusCode)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(1, client.addUser.count)

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
	client := &mockAddUserClient{}
	client.myDetails.roles = []string{"System Admin"}
	client.addUser.err = expectedErr
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := addUser(logger, client, template, "http://sirius")(w, r)
	assert.Equal(expectedErr, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(1, client.addUser.count)
	assert.Equal(0, template.count)
}

func TestPostAddUserMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	logger := log.New(ioutil.Discard, "", 0)
	client := &mockAddUserClient{}
	client.myDetails.err = expectedErr
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := addUser(logger, client, template, "http://sirius")(w, r)
	assert.Equal(expectedErr, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(0, client.addUser.count)
	assert.Equal(0, template.count)
}
