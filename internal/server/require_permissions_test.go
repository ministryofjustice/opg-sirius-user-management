package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockRequirePermissionsClient struct {
	count       int
	permissions sirius.PermissionSet
	err         error
}

func (m *mockRequirePermissionsClient) MyPermissions(ctx sirius.Context) (sirius.PermissionSet, error) {
	m.count++

	return m.permissions, m.err
}

func TestRequirePermissions(t *testing.T) {
	assert := assert.New(t)

	client := &mockRequirePermissionsClient{}
	client.permissions = sirius.PermissionSet{"team": sirius.PermissionGroup{Permissions: []string{"get"}}}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := requirePermissions(client, PermissionRequest{"team", http.MethodGet})(func(w http.ResponseWriter, r *http.Request) error {
		return StatusError(http.StatusOK)
	})(w, r)

	assert.Equal(StatusError(http.StatusOK), err)
	assert.Equal(1, client.count)
}

func TestRequirePermissionsChecksAll(t *testing.T) {
	assert := assert.New(t)

	client := &mockRequirePermissionsClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := requirePermissions(client, PermissionRequest{"team", http.MethodGet}, PermissionRequest{"team", http.MethodPost})(func(w http.ResponseWriter, r *http.Request) error {
		return StatusError(http.StatusOK)
	})

	client.permissions = sirius.PermissionSet{"team": sirius.PermissionGroup{Permissions: []string{http.MethodGet, http.MethodPost}}}
	err := handler(w, r)
	assert.Equal(StatusError(http.StatusOK), err)
	assert.Equal(1, client.count)

	client.permissions = sirius.PermissionSet{"team": sirius.PermissionGroup{Permissions: []string{http.MethodGet, http.MethodPost, http.MethodPatch}}}
	err = handler(w, r)
	assert.Equal(StatusError(http.StatusOK), err)
	assert.Equal(2, client.count)

	client.permissions = sirius.PermissionSet{"team": sirius.PermissionGroup{Permissions: []string{http.MethodGet}}}
	err = handler(w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)
	assert.Equal(3, client.count)

	client.permissions = sirius.PermissionSet{"user": sirius.PermissionGroup{Permissions: []string{http.MethodGet, http.MethodPost}}}
	err = handler(w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)
	assert.Equal(4, client.count)
}

func TestRequirePermissionsMissingPermission(t *testing.T) {
	assert := assert.New(t)

	client := &mockRequirePermissionsClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := requirePermissions(client, PermissionRequest{"team", http.MethodGet})(func(w http.ResponseWriter, r *http.Request) error {
		return StatusError(http.StatusOK)
	})(w, r)

	assert.Equal(StatusError(http.StatusForbidden), err)
	assert.Equal(1, client.count)
}

func TestRequirePermissionsError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockRequirePermissionsClient{}
	client.err = expectedErr

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := requirePermissions(client, PermissionRequest{"team", http.MethodGet})(func(w http.ResponseWriter, r *http.Request) error {
		return StatusError(http.StatusOK)
	})(w, r)

	assert.Equal(expectedErr, err)
	assert.Equal(1, client.count)
}
