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

type mockResendConfirmationClient struct {
	count     int
	lastCtx   sirius.Context
	lastEmail string
	err       error
}

func (m *mockResendConfirmationClient) ResendConfirmation(ctx sirius.Context, email string) error {
	m.count += 1
	m.lastCtx = ctx
	m.lastEmail = email

	return m.err
}

func (m *mockResendConfirmationClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-users": sirius.PermissionGroup{Permissions: []string{"put"}}}
}

func TestGetResendConfirmation(t *testing.T) {
	assert := assert.New(t)

	client := &mockResendConfirmationClient{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := resendConfirmation(nil, nil)(w, r)
	assert.Equal(handler.Redirect("/users"), err)
}

func TestGetResendConfirmationNoPermission(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, sirius.PermissionSet{}))

	err := resendConfirmation(nil, nil)(w, r)
	assert.Equal(handler.Status(http.StatusForbidden), err)
}

func TestPostResendConfirmation(t *testing.T) {
	assert := assert.New(t)

	client := &mockResendConfirmationClient{}
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("email=a&id=b"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := resendConfirmation(client, template.Func())(w, r)
	assert.Nil(err)

	assert.Equal(1, client.count)
	assert.Equal(getContext(r), client.lastCtx)
	assert.Equal("a", client.lastEmail)

	assert.Equal(1, template.count)
	assert.Equal(resendConfirmationVars{
		Path:  "/path",
		ID:    "b",
		Email: "a",
	}, template.lastVars)
}

func TestPostResendConfirmationError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockResendConfirmationClient{
		err: expectedErr,
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := resendConfirmation(client, nil)(w, r)
	assert.Equal(expectedErr, err)
}
