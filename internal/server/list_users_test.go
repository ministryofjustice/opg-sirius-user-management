package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/stretchr/testify/assert"
)

type mockListUsersClient struct {
	count      int
	lastCtx    sirius.Context
	lastSearch string
	err        error
	data       []sirius.User
}

func (m *mockListUsersClient) SearchUsers(ctx sirius.Context, search string) ([]sirius.User, error) {
	m.count += 1
	m.lastCtx = ctx
	m.lastSearch = search

	return m.data, m.err
}

func (m *mockListUsersClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-users": sirius.PermissionGroup{Permissions: []string{"put"}}}
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
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?search=milo", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	handler := listUsers(client, template.Func())
	err := handler(w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(1, client.count)
	assert.Equal("milo", client.lastSearch)

	assert.Equal(1, template.count)
	assert.Equal(listUsersVars{
		Path:   "/path",
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

func TestListUsersNoPermission(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, sirius.PermissionSet{}))

	err := listUsers(nil, nil)(w, r)
	assert.Equal(handler.Status(http.StatusForbidden), err)
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
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	handler := listUsers(client, template.Func())
	err := handler(w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(0, client.count)

	assert.Equal(1, template.count)
	assert.Equal(listUsersVars{
		Path:   "/path",
		Search: "",
		Users:  nil,
	}, template.lastVars)
}

func TestListUsersClientError(t *testing.T) {
	assert := assert.New(t)

	client := &mockListUsersClient{
		err: sirius.ClientError("problem"),
	}
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?search=m", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	handler := listUsers(client, template.Func())
	err := handler(w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal(listUsersVars{
		Path:   "/path",
		Search: "m",
		Users:  nil,
		Errors: sirius.ValidationErrors{
			"search": {
				"": "problem",
			},
		},
	}, template.lastVars)
}

func TestListUsersSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("err")
	client := &mockListUsersClient{err: expectedErr}
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/?search=long", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	handler := listUsers(client, template.Func())
	err := handler(w, r)

	assert.Equal(expectedErr, err)
	assert.Equal(0, template.count)
}

func TestPostListUsers(t *testing.T) {
	assert := assert.New(t)
	template := &mockTemplateFn{}

	client := &mockListUsersClient{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := listUsers(nil, template.Func())(w, r)

	assert.Equal(handler.Status(http.StatusMethodNotAllowed), err)

	assert.Equal(0, template.count)
}
