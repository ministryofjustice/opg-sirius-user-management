package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockListUsersClient struct {
	count          int
	myDetailsCount int
	lastCookies    []*http.Cookie
	lastSearch     string
	err            error
	data           []sirius.User
	myDetails      sirius.MyDetails
}

func (m *mockListUsersClient) SearchUsers(ctx context.Context, cookies []*http.Cookie, search string) ([]sirius.User, error) {
	m.count += 1
	m.lastCookies = cookies
	m.lastSearch = search

	return m.data, m.err
}

func (m *mockListUsersClient) MyDetails(ctx context.Context, cookies []*http.Cookie) (sirius.MyDetails, error) {
	m.myDetailsCount += 1
	m.lastCookies = cookies

	return m.myDetails, m.err
}

func TestListUsers(t *testing.T) {
	assert := assert.New(t)

	data := []sirius.User{
		{
			ID:          29,
			DisplayName: "Milo Nihei",
			Email:       "milo.nihei@opgtest.com",
			Status:      "Active",
		},
	}
	client := &mockListUsersClient{
		myDetails: sirius.MyDetails{
			Roles: []string{"System Admin"},
		},
		data: data,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?search=milo", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	handler := listUsers(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(r.Cookies(), client.lastCookies)

	assert.Equal(1, client.myDetailsCount)
	assert.Equal(1, client.count)
	assert.Equal("milo", client.lastSearch)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(listUsersVars{
		Path:      "/path",
		SiriusURL: "http://sirius",

		Search: "milo",
		Users: []sirius.User{
			{
				ID:          29,
				DisplayName: "Milo Nihei",
				Email:       "milo.nihei@opgtest.com",
				Status:      "Active",
			},
		},
	}, template.lastVars)
}

func TestListUsersRequiresSearch(t *testing.T) {
	assert := assert.New(t)

	data := []sirius.User{
		{
			ID:          29,
			DisplayName: "Milo Nihei",
			Email:       "milo.nihei@opgtest.com",
			Status:      "Active",
		},
	}
	client := &mockListUsersClient{
		myDetails: sirius.MyDetails{
			Roles: []string{"System Admin"},
		},
		data: data,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	handler := listUsers(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(r.Cookies(), client.lastCookies)

	assert.Equal(1, client.myDetailsCount)
	assert.Equal(0, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(listUsersVars{
		Path:      "/path",
		SiriusURL: "http://sirius",

		Search: "",
		Users:  nil,
	}, template.lastVars)
}

func TestListUsersWarnsShortSearches(t *testing.T) {
	assert := assert.New(t)

	data := []sirius.User{
		{
			ID:          29,
			DisplayName: "Milo Nihei",
			Email:       "milo.nihei@opgtest.com",
			Status:      "Active",
		},
	}
	client := &mockListUsersClient{
		myDetails: sirius.MyDetails{
			Roles: []string{"System Admin"},
		},
		data: data,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?search=m", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	handler := listUsers(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(r.Cookies(), client.lastCookies)

	assert.Equal(1, client.myDetailsCount)
	assert.Equal(0, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(listUsersVars{
		Path:      "/path",
		SiriusURL: "http://sirius",

		Search: "m",
		Users:  nil,
		Errors: sirius.ValidationErrors{
			"search": {
				"": "Search term must be at least three characters",
			},
		},
	}, template.lastVars)
}

func TestListUsersUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockListUsersClient{err: sirius.ErrUnauthorized}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := listUsers(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Equal(sirius.ErrUnauthorized, err)

	assert.Equal(0, template.count)
}

func TestListUsersMissingRole(t *testing.T) {
	assert := assert.New(t)

	client := &mockListUsersClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := listUsers(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Equal(StatusError(http.StatusForbidden), err)

	assert.Equal(0, template.count)
}

func TestListUsersSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockListUsersClient{err: errors.New("err")}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := listUsers(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Equal("err", err.Error())

	assert.Equal(0, template.count)
}

func TestPostListUsers(t *testing.T) {
	assert := assert.New(t)
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)

	handler := listUsers(nil, nil, template, "http://sirius")
	err := handler(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, template.count)
}
