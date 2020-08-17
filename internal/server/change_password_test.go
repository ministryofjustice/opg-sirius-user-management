package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockChangePasswordClient struct {
	errors            bool
	count             int
	old, new, confirm string
}

func (m *mockChangePasswordClient) ChangePassword(ctx context.Context, old, new, confirm string) error {
	m.count += 1
	m.old = old
	m.new = new
	m.confirm = confirm

	if m.errors {
		return errors.New("err")
	}

	return nil
}

func TestGetChangePassword(t *testing.T) {
	testCases := map[string]struct {
		url      string
		hasError bool
	}{
		"Default": {
			url:      "/",
			hasError: false,
		},
		"WithError": {
			url:      "/?error=1",
			hasError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			templates := &mockTemplates{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", tc.url, nil)

			changePassword(nil, templates).ServeHTTP(w, r)

			resp := w.Result()
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			assert.Equal(t, 1, templates.count)
			assert.Equal(t, "change-password.gotmpl", templates.lastName)
			assert.Equal(t, changePasswordVars{HasError: tc.hasError}, templates.lastVars)
		})
	}
}

func TestPostChangePassword(t *testing.T) {
	testCases := map[string]struct {
		errors     bool
		redirectTo string
	}{
		"Default": {
			errors:     false,
			redirectTo: "/my-details",
		},
		"WithError": {
			errors:     true,
			redirectTo: "/change-password?error=1",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &mockChangePasswordClient{errors: tc.errors}

			form := url.Values{
				"old-password":         {"old"},
				"new-password":         {"new"},
				"new-password-confirm": {"new-2"},
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "", strings.NewReader(form.Encode()))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			changePassword(client, nil).ServeHTTP(w, r)

			resp := w.Result()
			assert.Equal(t, http.StatusFound, resp.StatusCode)
			assert.Equal(t, tc.redirectTo, resp.Header.Get("Location"))

			assert.Equal(t, 1, client.count)
			assert.Equal(t, "old", client.old)
			assert.Equal(t, "new", client.new)
			assert.Equal(t, "new-2", client.confirm)
		})
	}
}

func TestOptionsChangePassword(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("OPTIONS", "", nil)

	changePassword(nil, nil).ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}
