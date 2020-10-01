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
	err            error
	data           []sirius.User
	myDetails      sirius.MyDetails
}

func (m *mockListUsersClient) ListUsers(ctx context.Context, cookies []*http.Cookie) ([]sirius.User, error) {
	m.count += 1
	m.lastCookies = cookies

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
	r, _ := http.NewRequest("GET", "/path", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	handler := listUsers(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(r.Cookies(), client.lastCookies)

	assert.Equal(1, client.myDetailsCount)
	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(listUsersVars{
		Path:      "/path",
		SiriusURL: "http://sirius",

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

func TestListUsersSearchesByNameOrEmail(t *testing.T) {
	assert := assert.New(t)

	data := []sirius.User{
		{
			ID:          29,
			DisplayName: "Milo Nihei",
			Email:       "milo.nihei@opgtest.com",
			Status:      "Active",
		},
		{
			ID:          34,
			DisplayName: "Blair Lemmons",
			Email:       "blair@opgtest.com",
			Status:      "Active",
		},
		{
			ID:          83,
			DisplayName: "Lori Rajtar",
			Email:       "Lori.lemmons@opgtest.com",
			Status:      "Locked",
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
	r, _ := http.NewRequest("GET", "/path?search=LEMMONS", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	handler := listUsers(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Nil(err)

	assert.Equal(listUsersVars{
		Path:      "/path",
		SiriusURL: "http://sirius",

		Search: "LEMMONS",
		Users: []sirius.User{
			{
				ID:          34,
				DisplayName: "Blair Lemmons",
				Email:       "blair@opgtest.com",
				Status:      "Active",
			},
			{
				ID:          83,
				DisplayName: "Lori Rajtar",
				Email:       "Lori.lemmons@opgtest.com",
				Status:      "Locked",
			},
		},
	}, template.lastVars)
}

func TestListUsersSearchesByStatus(t *testing.T) {
	assert := assert.New(t)

	data := []sirius.User{
		{
			ID:          29,
			DisplayName: "Milo Nihei",
			Email:       "milo.nihei@opgtest.com",
			Status:      "Active",
		},
		{
			ID:          83,
			DisplayName: "Lori Rajtar",
			Email:       "Lori.lemmons@opgtest.com",
			Status:      "Locked",
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
	r, _ := http.NewRequest("GET", "/path?search=locked", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	handler := listUsers(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Nil(err)

	assert.Equal(listUsersVars{
		Path:      "/path",
		SiriusURL: "http://sirius",

		Search: "locked",
		Users: []sirius.User{
			{
				ID:          83,
				DisplayName: "Lori Rajtar",
				Email:       "Lori.lemmons@opgtest.com",
				Status:      "Locked",
			},
		},
	}, template.lastVars)
}
