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

type mockDeleteUserClient struct {
	user struct {
		count       int
		lastCookies []*http.Cookie
		lastID      int
		data        sirius.AuthUser
		err         error
	}

	deleteUser struct {
		count       int
		lastCookies []*http.Cookie
		lastUserID  int
		err         error
	}
}

func (m *mockDeleteUserClient) User(ctx context.Context, cookies []*http.Cookie, id int) (sirius.AuthUser, error) {
	m.user.count += 1
	m.user.lastCookies = cookies
	m.user.lastID = id

	return m.user.data, m.user.err
}

func (m *mockDeleteUserClient) DeleteUser(ctx context.Context, cookies []*http.Cookie, userID int) error {
	m.deleteUser.count += 1
	m.deleteUser.lastCookies = cookies
	m.deleteUser.lastUserID = userID

	return m.deleteUser.err
}

func TestGetDeleteUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeleteUserClient{}
	client.user.data = sirius.AuthUser{Firstname: "test"}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/delete-user/123", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := deleteUser(client, template, "http://sirius")(w, r)
	assert.Nil(err)

	assert.Equal(1, client.user.count)
	assert.Equal(123, client.user.lastID)
	assert.Equal(0, client.deleteUser.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(deleteUserVars{
		Path:      "/delete-user/123",
		SiriusURL: "http://sirius",
		User:      client.user.data,
	}, template.lastVars)
}

func TestGetDeleteUserError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockDeleteUserClient{}
	client.user.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/delete-user/123", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := deleteUser(client, template, "http://sirius")(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.user.count)
	assert.Equal(123, client.user.lastID)
	assert.Equal(0, client.deleteUser.count)
}

func TestGetDeleteUserBadPath(t *testing.T) {
	for name, path := range map[string]string{
		"empty":       "/delete-user/",
		"non-numeric": "/delete-user/hello",
		"suffixed":    "/delete-user/123/no",
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockDeleteUserClient{}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", path, nil)
			r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

			err := deleteUser(client, template, "http://sirius")(w, r)
			assert.Equal(StatusError(http.StatusNotFound), err)

			assert.Equal(0, client.user.count)
			assert.Equal(0, client.deleteUser.count)
			assert.Equal(0, template.count)
		})
	}
}

func TestPostDeleteUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeleteUserClient{}
	client.user.data = sirius.AuthUser{Firstname: "test"}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/delete-user/123", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := deleteUser(client, template, "http://sirius")(w, r)
	assert.Equal(RedirectError("/users"), err)

	assert.Equal(1, client.deleteUser.count)
	assert.Equal(r.Cookies(), client.deleteUser.lastCookies)
	assert.Equal(123, client.deleteUser.lastUserID)

	assert.Equal(0, client.user.count)
	assert.Equal(0, template.count)
}

func TestPostDeleteUserError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockDeleteUserClient{}
	client.deleteUser.err = expectedErr
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/delete-user/123", nil)

	err := deleteUser(client, template, "http://sirius")(w, r)
	assert.Equal(expectedErr, err)

	assert.Equal(1, client.deleteUser.count)
	assert.Equal(0, client.user.count)
	assert.Equal(0, template.count)
}

func TestPutDeleteUser(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/delete-user/123", nil)

	err := deleteUser(nil, nil, "http://sirius")(w, r)
	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
}
