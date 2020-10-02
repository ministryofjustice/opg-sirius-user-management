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
	count       int
	lastCookies []*http.Cookie
	lastSearch  string
	err         error
	data        []sirius.User
}

func (m *mockListUsersClient) SearchUsers(ctx context.Context, cookies []*http.Cookie, search string) ([]sirius.User, error) {
	m.count += 1
	m.lastCookies = cookies
	m.lastSearch = search

	return m.data, m.err
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

func TestListUsersSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("err")
	client := &mockListUsersClient{err: expectedErr}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/?search=long", nil)

	handler := listUsers(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Equal(expectedErr, err)
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
