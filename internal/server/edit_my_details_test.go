package server

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

func TestGetEditMyDetails(t *testing.T) {
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

	editMyDetails(nil, client, template, "http://sirius").ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(r.Cookies(), client.lastCookies)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editMyDetailsVars{
		Path:        "/path",
		SiriusURL:   "http://sirius",
		PhoneNumber: "123",
	}, template.lastVars)
}

func TestGetEditMyDetailsUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockMyDetailsClient{err: sirius.ErrUnauthorized}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	editMyDetails(nil, client, template, "http://sirius").ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusFound, resp.StatusCode)
	assert.Equal("http://sirius/auth", resp.Header.Get("Location"))

	assert.Equal(0, template.count)
}

func TestGetEditMyDetailsSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	logger := log.New(ioutil.Discard, "", 0)
	client := &mockMyDetailsClient{err: errors.New("err")}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	editMyDetails(logger, client, template, "http://sirius").ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(0, template.count)
}
