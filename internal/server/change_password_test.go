package server

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockChangePasswordClient struct {
	count                  int
	lastCookies            []*http.Cookie
	lastExistingPassword   string
	lastNewPassword        string
	lastNewPasswordConfirm string
	err                    error
}

func (m *mockChangePasswordClient) ChangePassword(ctx context.Context, cookies []*http.Cookie, existingPassword, newPassword, newPasswordConfirm string) error {
	m.count += 1
	m.lastCookies = cookies
	m.lastExistingPassword = existingPassword
	m.lastNewPassword = newPassword
	m.lastNewPasswordConfirm = newPasswordConfirm

	return m.err
}

func TestGetChangePassword(t *testing.T) {
	assert := assert.New(t)

	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	changePassword(nil, nil, template, "/prefix", "http://sirius").ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(changePasswordVars{
		Path:      "/path",
		SiriusURL: "http://sirius",
		Prefix:    "/prefix",
	}, template.lastVars)
}

func TestGetChangePasswordWithError(t *testing.T) {
	assert := assert.New(t)

	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?error=Something+happened", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	changePassword(nil, nil, template, "/prefix", "http://sirius").ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(changePasswordVars{
		Path:      "/path",
		SiriusURL: "http://sirius",
		Prefix:    "/prefix",
		Error:     "Something happened",
	}, template.lastVars)
}

func TestPostChangePassword(t *testing.T) {
	assert := assert.New(t)

	client := &mockChangePasswordClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("currentpassword=a&password1=b&password2=c"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	changePassword(nil, client, template, "/prefix", "http://sirius").ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusFound, resp.StatusCode)
	assert.Equal("/prefix/my-details", resp.Header.Get("Location"))
	assert.Equal(r.Cookies(), client.lastCookies)
	assert.Equal("a", client.lastExistingPassword)
	assert.Equal("b", client.lastNewPassword)
	assert.Equal("c", client.lastNewPasswordConfirm)

	assert.Equal(0, template.count)
}

func TestPostChangePasswordUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockChangePasswordClient{err: sirius.ErrUnauthorized}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	changePassword(nil, client, template, "/prefix", "http://sirius").ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusFound, resp.StatusCode)
	assert.Equal("http://sirius/auth", resp.Header.Get("Location"))

	assert.Equal(0, template.count)
}

func TestPostChangePasswordSiriusError(t *testing.T) {
	assert := assert.New(t)

	client := &mockChangePasswordClient{err: sirius.ClientError("Something happened")}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	changePassword(nil, client, template, "/prefix", "http://sirius").ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusFound, resp.StatusCode)
	assert.Equal("/prefix/change-password?error=Something+happened", resp.Header.Get("Location"))

	assert.Equal(0, template.count)
}

func TestPostChangePasswordOtherError(t *testing.T) {
	assert := assert.New(t)

	logger := log.New(ioutil.Discard, "", 0)
	client := &mockChangePasswordClient{err: errors.New("oops")}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	changePassword(logger, client, template, "/prefix", "http://sirius").ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusFound, resp.StatusCode)
	assert.Equal("/prefix/change-password", resp.Header.Get("Location"))

	assert.Equal(0, template.count)
}
