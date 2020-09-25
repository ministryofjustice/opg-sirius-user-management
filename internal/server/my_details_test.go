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
	count       int
	lastCookies []*http.Cookie
	err         error
	data        sirius.MyDetails
	errors      sirius.ValidationErrors
}

func (m *mockMyDetailsClient) MyDetails(ctx context.Context, cookies []*http.Cookie) (sirius.MyDetails, error) {
	m.count += 1
	m.lastCookies = cookies

	return m.data, m.err
}

func (m *mockMyDetailsClient) EditMyDetails(ctx context.Context, cookies []*http.Cookie, id int, phoneNumber string) (sirius.ValidationErrors, error) {
	return nil, nil
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

	myDetails(nil, client, template, "http://sirius").ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(r.Cookies(), client.lastCookies)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(myDetailsVars{
		Path:         "/path",
		SiriusURL:    "http://sirius",
		ID:           123,
		Firstname:    "John",
		Surname:      "Doe",
		Email:        "john@doe.com",
		PhoneNumber:  "123",
		Organisation: "COP User",
		Roles:        []string{"A", "B"},
		Teams:        []string{"A Team"},
	}, template.lastVars)
}

func TestGetMyDetailsUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockMyDetailsClient{err: sirius.ErrUnauthorized}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	myDetails(nil, client, template, "http://sirius").ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusFound, resp.StatusCode)
	assert.Equal("http://sirius/auth", resp.Header.Get("Location"))

	assert.Equal(0, template.count)
}

func TestGetMyDetailsSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	logger := log.New(ioutil.Discard, "", 0)
	client := &mockMyDetailsClient{err: errors.New("err")}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	myDetails(logger, client, template, "http://sirius").ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(0, template.count)
}

func TestPostMyDetails(t *testing.T) {
	assert := assert.New(t)
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)

	myDetails(nil, nil, template, "http://sirius").ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusMethodNotAllowed, resp.StatusCode)
	assert.Equal(0, template.count)
}
