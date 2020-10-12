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

type mockAllowRolesClient struct {
	count int
	roles []string
	err   error
}

func (m *mockAllowRolesClient) MyDetails(ctx context.Context, cookies []*http.Cookie) (sirius.MyDetails, error) {
	m.count += 1

	return sirius.MyDetails{Roles: m.roles}, m.err
}

func TestAllowRoles(t *testing.T) {
	assert := assert.New(t)

	client := &mockAllowRolesClient{}
	client.roles = []string{"System Admin"}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := allowRoles(client, "System Admin")(func(w http.ResponseWriter, r *http.Request) error {
		return StatusError(http.StatusTeapot)
	})(w, r)

	assert.Equal(StatusError(http.StatusTeapot), err)
	assert.Equal(1, client.count)
}

func TestAllowRolesMultipleChoices(t *testing.T) {
	assert := assert.New(t)

	client := &mockAllowRolesClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	handler := allowRoles(client, "System Admin", "Manager")(func(w http.ResponseWriter, r *http.Request) error {
		return StatusError(http.StatusTeapot)
	})

	client.roles = []string{"System Admin", "Manager"}
	err := handler(w, r)
	assert.Equal(StatusError(http.StatusTeapot), err)
	assert.Equal(1, client.count)

	client.roles = []string{"System Admin"}
	err = handler(w, r)
	assert.Equal(StatusError(http.StatusTeapot), err)
	assert.Equal(2, client.count)

	client.roles = []string{"Manager"}
	err = handler(w, r)
	assert.Equal(StatusError(http.StatusTeapot), err)
	assert.Equal(3, client.count)

	client.roles = []string{"What"}
	err = handler(w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)
	assert.Equal(4, client.count)
}

func TestAllowRolesMissingRole(t *testing.T) {
	assert := assert.New(t)

	client := &mockAllowRolesClient{}
	client.roles = []string{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := allowRoles(client, "System Admin")(func(w http.ResponseWriter, r *http.Request) error {
		return StatusError(http.StatusTeapot)
	})(w, r)

	assert.Equal(StatusError(http.StatusForbidden), err)
	assert.Equal(1, client.count)
}

func TestAllowRolesMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockAllowRolesClient{}
	client.err = expectedErr

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := allowRoles(client, "System Admin")(func(w http.ResponseWriter, r *http.Request) error {
		return StatusError(http.StatusTeapot)
	})(w, r)

	assert.Equal(expectedErr, err)
	assert.Equal(1, client.count)
}
