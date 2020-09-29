package server

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockMyDetailsClient struct {
	count            int
	permissionsCount int
	lastCookies      []*http.Cookie
	err              error
	permissionsErr   error
	data             sirius.MyDetails
	hasPermission    bool
	lastPermission   struct {
		Group  string
		Method string
	}
}

func (m *mockMyDetailsClient) MyDetails(ctx context.Context, cookies []*http.Cookie) (sirius.MyDetails, error) {
	m.count += 1
	m.lastCookies = cookies

	return m.data, m.err
}

func (m *mockMyDetailsClient) HasPermission(ctx context.Context, cookies []*http.Cookie, group string, method string) (bool, error) {
	m.permissionsCount += 1
	m.lastCookies = cookies
	m.lastPermission.Group = group
	m.lastPermission.Method = method

	return m.hasPermission, m.permissionsErr
}

func TestGetMyDetails(t *testing.T) {
	assert := assert.New(t)

	data := sirius.MyDetails{
		ID:          123,
		Firstname:   "John",
		Surname:     "Doe",
		Email:       "john@doe.com",
		PhoneNumber: "123",
		Roles:       []string{"A", "COP User", "B"},
		Teams: []sirius.MyDetailsTeam{
			{DisplayName: "A Team"},
		},
	}
	client := &mockMyDetailsClient{data: data}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	handler := myDetails(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(r.Cookies(), client.lastCookies)

	assert.Equal(1, client.count)
	assert.Equal(1, client.permissionsCount)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(myDetailsVars{
		Path:               "/path",
		SiriusURL:          "http://sirius",
		ID:                 123,
		Firstname:          "John",
		Surname:            "Doe",
		Email:              "john@doe.com",
		PhoneNumber:        "123",
		Organisation:       "COP User",
		Roles:              []string{"A", "B"},
		Teams:              []string{"A Team"},
		CanEditPhoneNumber: false,
	}, template.lastVars)
}

func TestGetMyDetailsUsesPermission(t *testing.T) {
	assert := assert.New(t)

	data := sirius.MyDetails{
		ID:          123,
		Firstname:   "John",
		Surname:     "Doe",
		Email:       "john@doe.com",
		PhoneNumber: "123",
		Roles:       []string{"A", "COP User", "B"},
		Teams: []sirius.MyDetailsTeam{
			{DisplayName: "A Team"},
		},
	}
	client := &mockMyDetailsClient{data: data, hasPermission: true}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	handler := myDetails(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(r.Cookies(), client.lastCookies)
	assert.Equal("user", client.lastPermission.Group)
	assert.Equal("patch", client.lastPermission.Method)

	assert.Equal(1, client.count)
	assert.Equal(1, client.permissionsCount)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(myDetailsVars{
		Path:               "/path",
		SiriusURL:          "http://sirius",
		ID:                 123,
		Firstname:          "John",
		Surname:            "Doe",
		Email:              "john@doe.com",
		PhoneNumber:        "123",
		Organisation:       "COP User",
		Roles:              []string{"A", "B"},
		Teams:              []string{"A Team"},
		CanEditPhoneNumber: true,
	}, template.lastVars)
}

func TestGetMyDetailsUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockMyDetailsClient{err: sirius.ErrUnauthorized}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := myDetails(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Equal(sirius.ErrUnauthorized, err)

	assert.Equal(0, template.count)
}

func TestGetMyDetailsSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	logger := log.New(ioutil.Discard, "", 0)
	client := &mockMyDetailsClient{err: errors.New("err")}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := myDetails(logger, client, template, "http://sirius")
	err := handler(w, r)

	status, ok := err.(StatusError)
	assert.True(ok)
	assert.Equal(http.StatusInternalServerError, status.Code())

	assert.Equal(0, template.count)
}

func TestPostMyDetails(t *testing.T) {
	assert := assert.New(t)
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)

	handler := myDetails(nil, nil, template, "http://sirius")
	err := handler(w, r)

	status, ok := err.(StatusError)
	assert.True(ok)
	assert.Equal(http.StatusMethodNotAllowed, status.Code())

	assert.Equal(0, template.count)
}
