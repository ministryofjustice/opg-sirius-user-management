package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	assert := assert.New(t)

	nilHandler := func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}

	handler := New("http://redirect", "http://auth", func(w http.ResponseWriter, r *http.Request, code int, err error) {
		assert.Fail("an error should not be handled")
	})(nilHandler)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestHandlerRedirect(t *testing.T) {
	assert := assert.New(t)

	redirectHandler := func(w http.ResponseWriter, r *http.Request) error {
		return Redirect("/x/y/z")
	}

	handler := New("http://redirect", "http://auth", func(w http.ResponseWriter, r *http.Request, code int, err error) {
		assert.Fail("an error should not be handled")
	})(redirectHandler)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusFound, resp.StatusCode)
	assert.Equal("http://redirect/x/y/z", resp.Header.Get("Location"))
}

func TestHandlerForbidden(t *testing.T) {
	assert := assert.New(t)
	errCh := make(chan error, 1)

	statusHandler := func(w http.ResponseWriter, r *http.Request) error {
		return Status(http.StatusForbidden)
	}

	handler := New("http://redirect", "http://auth", func(w http.ResponseWriter, r *http.Request, code int, err error) {
		errCh <- err
	})(statusHandler)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusForbidden, resp.StatusCode)

	select {
	case err := <-errCh:
		assert.Equal("403 Forbidden", err.Error())
	case <-time.After(time.Millisecond):
		assert.Fail("an error should have been handled")
	}
}

func TestHandlerNotFound(t *testing.T) {
	assert := assert.New(t)
	errCh := make(chan error, 1)

	statusHandler := func(w http.ResponseWriter, r *http.Request) error {
		return Status(http.StatusNotFound)
	}

	handler := New("http://redirect", "http://auth", func(w http.ResponseWriter, r *http.Request, code int, err error) {
		errCh <- err
	})(statusHandler)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusNotFound, resp.StatusCode)

	select {
	case err := <-errCh:
		assert.Equal("404 Not Found", err.Error())
	case <-time.After(time.Millisecond):
		assert.Fail("an error should have been handled")
	}
}

func TestHandlerStatus(t *testing.T) {
	assert := assert.New(t)
	errCh := make(chan error, 1)

	statusHandler := func(w http.ResponseWriter, r *http.Request) error {
		return Status(http.StatusTeapot)
	}

	handler := New("http://redirect", "http://auth", func(w http.ResponseWriter, r *http.Request, code int, err error) {
		errCh <- err
	})(statusHandler)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusInternalServerError, resp.StatusCode)

	select {
	case err := <-errCh:
		assert.Equal("418 I'm a teapot", err.Error())
	case <-time.After(time.Millisecond):
		assert.Fail("an error should have been handled")
	}
}

func TestHandlerUnauthorizedStatus(t *testing.T) {
	assert := assert.New(t)
	errCh := make(chan error, 1)

	statusHandler := func(w http.ResponseWriter, r *http.Request) error {
		return Status(http.StatusUnauthorized)
	}

	handler := New("http://redirect", "http://auth", func(w http.ResponseWriter, r *http.Request, code int, err error) {
		errCh <- err
	})(statusHandler)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusFound, resp.StatusCode)
	assert.Equal("http://auth", resp.Header.Get("Location"))
}
