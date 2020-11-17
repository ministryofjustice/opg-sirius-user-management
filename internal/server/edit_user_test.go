package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
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

func TestGetEditUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditUserClient{}
	client.user.data = sirius.AuthUser{Firstname: "test"}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/edit-user/123", nil)

	err := editUser(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.user.count)
	assert.Equal(123, client.user.lastID)
	assert.Equal(0, client.editUser.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editUserVars{
		Path: "/edit-user/123",
		User: client.user.data,
	}, template.lastVars)
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
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", path, nil)

			err := editUser(client, template)(w, r)
			assert.Equal(StatusError(http.StatusNotFound), err)

			assert.Equal(0, client.user.count)
			assert.Equal(0, client.editUser.count)
			assert.Equal(0, template.count)
		})
	}
}

func TestPostEditUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditUserClient{}
	client.user.data = sirius.AuthUser{Firstname: "test"}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/edit-user/123", strings.NewReader("email=a&firstname=b&surname=c&organisation=d&roles=e&roles=f&locked=Yes&suspended=No"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := editUser(client, template)(w, r)
	assert.Nil(err)

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
	assert.Equal("page", template.lastName)
	assert.Equal(editUserVars{
		Path:    "/edit-user/123",
		Success: true,
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
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/edit-user/123", strings.NewReader("email=a&firstname=b&surname=c&organisation=d&roles=e&roles=f&locked=Yes&suspended=No"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := editUser(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.editUser.count)
	assert.Equal(0, client.user.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editUserVars{
		Path: "/edit-user/123",
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
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/edit-user/123", nil)

	err := editUser(client, template)(w, r)
	assert.Equal(expectedErr, err)

	assert.Equal(1, client.editUser.count)
	assert.Equal(0, client.user.count)
	assert.Equal(0, template.count)
}
