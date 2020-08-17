package sirius

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type request struct {
	method, path string
	form         url.Values
}

func TestChangePassword(t *testing.T) {
	requests := make(chan request, 1)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		requests <- request{
			method: r.Method,
			path:   r.URL.Path,
			form:   r.PostForm,
		}
	}))
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.ChangePassword(context.Background(), "existing", "new", "new-2")
	assert.Nil(t, err)

	select {
	case r := <-requests:
		assert.Equal(t, http.MethodPost, r.method)
		assert.Equal(t, "/auth/change-password", r.path)
		assert.Equal(t, "existing", r.form.Get("existingPassword"))
		assert.Equal(t, "new", r.form.Get("password"))
		assert.Equal(t, "new-2", r.form.Get("confirmPassword"))

	case <-time.After(time.Millisecond):
		assert.Fail(t, "request did not happen in time")
	}
}

func TestChangePasswordWhenBadRequest(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusBadRequest)
	}))
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.ChangePassword(context.Background(), "existing", "new", "new-2")
	assert.NotNil(t, err)
}
