package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/stretchr/testify/assert"
)

type mockEditUserClient struct {
	user struct {
		count   int
		lastCtx sirius.Context
		lastID  int
		data    sirius.AuthUser
		err     error
	}

	editUser struct {
		count    int
		lastCtx  sirius.Context
		lastUser sirius.AuthUser
		err      error
	}

	roles struct {
		count   int
		lastCtx sirius.Context
		err     error
	}
}

func (m *mockEditUserClient) User(ctx sirius.Context, id int) (sirius.AuthUser, error) {
	m.user.count += 1
	m.user.lastCtx = ctx
	m.user.lastID = id

	return m.user.data, m.user.err
}

func (m *mockEditUserClient) EditUser(ctx sirius.Context, user sirius.AuthUser) error {
	m.editUser.count += 1
	m.editUser.lastCtx = ctx
	m.editUser.lastUser = user

	return m.editUser.err
}

func (m *mockEditUserClient) Roles(ctx sirius.Context) ([]string, error) {
	m.roles.count += 1
	m.roles.lastCtx = ctx

	return []string{"System Admin", "Manager"}, m.roles.err
}

func (m *mockEditUserClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-users": sirius.PermissionGroup{Permissions: []string{"put"}}}
}

func TestGetEditUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditUserClient{}
	client.user.data = sirius.AuthUser{Firstname: "test"}
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/edit-user/123", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := editUser(client, template.Func())(w, r)
	assert.Nil(err)

	assert.Equal(1, client.roles.count)
	assert.Equal(getContext(r), client.roles.lastCtx)
	assert.Equal(1, client.user.count)
	assert.Equal(123, client.user.lastID)
	assert.Equal(0, client.editUser.count)

	assert.Equal(1, template.count)
	assert.Equal(editUserVars{
		Path:  "/edit-user/123",
		User:  client.user.data,
		Roles: []string{"System Admin", "Manager"},
	}, template.lastVars)
}

func TestGetEditUserNoPermission(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, sirius.PermissionSet{}))

	err := editUser(nil, nil)(w, r)
	assert.Equal(handler.Status(http.StatusForbidden), err)
}

func TestGetEditUserBadPath(t *testing.T) {
	for name, path := range map[string]string{
		"empty":       "/edit-user/",
		"non-numeric": "/edit-user/hello",
		"suffixed":    "/edit-user/123/no",
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockEditUserClient{}
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", path, nil)
			r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

			err := editUser(nil, nil)(w, r)
			assert.Equal(handler.Status(http.StatusNotFound), err)
		})
	}
}

func TestPostEditUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditUserClient{}
	client.user.data = sirius.AuthUser{Firstname: "test"}
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/edit-user/123", strings.NewReader("email=a&firstname=b&surname=c&organisation=d&roles=e&roles=f&locked=Yes&suspended=No"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := editUser(client, template.Func())(w, r)
	assert.Nil(err)

	assert.Equal(1, client.roles.count)
	assert.Equal(getContext(r), client.roles.lastCtx)

	assert.Equal(1, client.editUser.count)
	assert.Equal(getContext(r), client.editUser.lastCtx)
	assert.Equal(sirius.AuthUser{
		ID:           123,
		Firstname:    "b",
		Surname:      "c",
		Organisation: "d",
		Roles:        []string{"e", "f"},
		Locked:       true,
		Suspended:    false,
	}, client.editUser.lastUser)

	assert.Equal(0, client.user.count)

	assert.Equal(1, template.count)
	assert.Equal(editUserVars{
		Path:    "/edit-user/123",
		Success: true,
		Roles:   []string{"System Admin", "Manager"},
		User: sirius.AuthUser{
			ID:           123,
			Email:        "a",
			Firstname:    "b",
			Surname:      "c",
			Organisation: "d",
			Roles:        []string{"e", "f"},
			Locked:       true,
			Suspended:    false,
		},
	}, template.lastVars)
}

func TestPostEditUserClientError(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditUserClient{}
	client.editUser.err = sirius.ClientError("something")
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/edit-user/123", strings.NewReader("email=a&firstname=b&surname=c&organisation=d&roles=e&roles=f&locked=Yes&suspended=No"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := editUser(client, template.Func())(w, r)
	assert.Nil(err)

	assert.Equal(1, client.roles.count)
	assert.Equal(1, client.editUser.count)
	assert.Equal(0, client.user.count)

	assert.Equal(1, template.count)
	assert.Equal(editUserVars{
		Path:  "/edit-user/123",
		Roles: []string{"System Admin", "Manager"},
		User: sirius.AuthUser{
			ID:           123,
			Firstname:    "b",
			Surname:      "c",
			Organisation: "d",
			Roles:        []string{"e", "f"},
			Locked:       true,
			Suspended:    false,
		},
		Errors: sirius.ValidationErrors{
			"firstname": {
				"": "something",
			},
		},
	}, template.lastVars)
}

func TestPostEditUserOtherError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockEditUserClient{}
	client.editUser.err = expectedErr
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/edit-user/123", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := editUser(client, template.Func())(w, r)
	assert.Equal(expectedErr, err)

	assert.Equal(1, client.roles.count)
	assert.Equal(1, client.editUser.count)
	assert.Equal(0, client.user.count)
	assert.Equal(0, template.count)
}

func TestPostEditUserRolesError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockEditUserClient{}
	client.roles.err = expectedErr
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/edit-user/123", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := editUser(client, template.Func())(w, r)
	assert.Equal(expectedErr, err)

	assert.Equal(1, client.roles.count)
	assert.Equal(0, client.editUser.count)
	assert.Equal(0, client.user.count)
	assert.Equal(0, template.count)
}
