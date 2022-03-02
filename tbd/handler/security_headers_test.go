package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	assert := assert.New(t)

	handler := SecurityHeaders(http.NotFoundHandler())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()

	assert.Equal("default-src 'self'", resp.Header.Get("Content-Security-Policy"))
	assert.Equal("same-origin", resp.Header.Get("Referrer-Policy"))
	assert.Equal("max-age=31536000; includeSubDomains; preload", resp.Header.Get("Strict-Transport-Security"))
	assert.Equal("nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal("SAMEORIGIN", resp.Header.Get("X-Frame-Options"))
	assert.Equal("1; mode=block", resp.Header.Get("X-XSS-Protection"))
}
