package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockTemplate struct {
	count    int
	lastName string
	lastVars interface{}
	err      error
}

func (m *mockTemplate) ExecuteTemplate(w io.Writer, name string, vars interface{}) error {
	m.count += 1
	m.lastName = name
	m.lastVars = vars
	return m.err
}

type mockErrorHandlerClient struct {
	count       int
	err         error
	permissions sirius.PermissionSet
}

func (m *mockErrorHandlerClient) MyPermissions(ctx sirius.Context) (sirius.PermissionSet, error) {
	m.count++

	return m.permissions, m.err
}

func contextWithLogger() (context.Context, *bytes.Buffer) {
	var buf bytes.Buffer
	logHandler := slog.NewJSONHandler(&buf, nil)
	ctx := telemetry.ContextWithLogger(context.Background(), slog.New(logHandler))

	return ctx, &buf
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*http.Handler)(nil), New(nil, nil, nil, "", "", ""))
}

func TestErrorHandler(t *testing.T) {
	assert := assert.New(t)

	client := &mockErrorHandlerClient{}
	tmplError := &mockTemplate{}

	wrap := errorHandler(client, tmplError, "/prefix", "http://sirius")
	handler := wrap(func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusTeapot)
		return nil
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()

	assert.Equal(1, client.count)

	assert.Equal(http.StatusTeapot, resp.StatusCode)
	assert.Equal(0, tmplError.count)
}

func TestErrorHandlerUnauthorized(t *testing.T) {
	assert := assert.New(t)

	client := &mockErrorHandlerClient{}
	tmplError := &mockTemplate{}

	wrap := errorHandler(client, tmplError, "/prefix", "http://sirius")
	handler := wrap(func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		return sirius.ErrUnauthorized
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusFound, resp.StatusCode)
	assert.Equal("http://sirius/auth?redirect=%2Fprefix%2Fpath", resp.Header.Get("Location"))

	assert.Equal(0, tmplError.count)
}

func TestErrorHandlerMyPermissionsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	ctx, logBuf := contextWithLogger()
	client := &mockErrorHandlerClient{}
	client.err = expectedError
	tmplError := &mockTemplate{}

	wrap := errorHandler(client, tmplError, "/prefix", "http://sirius")
	handler := wrap(func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		return sirius.ErrUnauthorized
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(ctx, "GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()

	data := map[string]string{}
	err := json.Unmarshal(logBuf.Bytes(), &data)
	assert.Nil(err)
	assert.Equal("oops", data["msg"])
	assert.Equal("ERROR", data["level"])

	assert.Equal(http.StatusInternalServerError, resp.StatusCode)

	assert.Equal(1, tmplError.count)
}

func TestErrorHandlerRedirect(t *testing.T) {
	assert := assert.New(t)

	client := &mockErrorHandlerClient{}
	tmplError := &mockTemplate{}

	wrap := errorHandler(client, tmplError, "/prefix", "http://sirius")
	handler := wrap(func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		return RedirectError("/here")
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusFound, resp.StatusCode)
	assert.Equal("/prefix/here", resp.Header.Get("Location"))

	assert.Equal(0, tmplError.count)
}

func TestErrorHandlerStatus(t *testing.T) {
	assert := assert.New(t)

	ctx, logBuf := contextWithLogger()
	client := &mockErrorHandlerClient{}
	tmplError := &mockTemplate{}

	wrap := errorHandler(client, tmplError, "/prefix", "http://sirius")
	handler := wrap(func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		return StatusError(http.StatusTeapot)
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(ctx, "GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusInternalServerError, resp.StatusCode)

	assert.Equal(1, tmplError.count)
	assert.Equal(errorVars{SiriusURL: "http://sirius", Code: http.StatusInternalServerError, Error: "418 I'm a teapot"}, tmplError.lastVars)

	data := map[string]string{}
	err := json.Unmarshal(logBuf.Bytes(), &data)
	assert.Nil(err)
	assert.Equal("418 I'm a teapot", data["msg"])
	assert.Equal("ERROR", data["level"])
}

func TestErrorHandlerStatusKnown(t *testing.T) {
	for name, code := range map[string]int{
		"Forbidden": http.StatusForbidden,
		"Not Found": http.StatusNotFound,
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			ctx, logBuf := contextWithLogger()
			client := &mockErrorHandlerClient{}
			tmplError := &mockTemplate{}

			wrap := errorHandler(client, tmplError, "/prefix", "http://sirius")
			handler := wrap(func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
				return StatusError(code)
			})

			w := httptest.NewRecorder()
			r, _ := http.NewRequestWithContext(ctx, "GET", "/path", nil)

			handler.ServeHTTP(w, r)

			resp := w.Result()
			assert.Equal(code, resp.StatusCode)

			assert.Equal(1, tmplError.count)
			assert.Equal(errorVars{SiriusURL: "http://sirius", Code: code, Error: fmt.Sprintf("%d %s", code, name)}, tmplError.lastVars)

			assert.Equal("", logBuf.String())
		})
	}
}

func TestErrorHandlerTemplateError(t *testing.T) {
	assert := assert.New(t)

	ctx, logBuf := contextWithLogger()
	client := &mockErrorHandlerClient{}
	tmplError := &mockTemplate{}
	tmplError.err = errors.New("could not render")

	wrap := errorHandler(client, tmplError, "/prefix", "http://sirius")
	handler := wrap(func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		return StatusError(http.StatusNotFound)
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(ctx, "GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusNotFound, resp.StatusCode)

	assert.Equal(1, tmplError.count)

	data := map[string]string{}
	err := json.Unmarshal(logBuf.Bytes(), &data)
	assert.Nil(err)
	assert.Equal("could not generate error template", data["msg"])
	assert.Equal("ERROR", data["level"])
	assert.Equal("could not render", data["err"])
}

func TestGetContext(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: "z3tVRZ00yx4dHz3KWYv3boLWHZ4/RsCsVAKbvo2SBNc%3D"})
	r.AddCookie(&http.Cookie{Name: "another", Value: "one"})

	ctx := getContext(r)
	assert.Equal(r.Context(), ctx.Context)
	assert.Equal(r.Cookies(), ctx.Cookies)
	assert.Equal("z3tVRZ00yx4dHz3KWYv3boLWHZ4/RsCsVAKbvo2SBNc=", ctx.XSRFToken)
}

func TestGetContextBadXSRFToken(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: "%"})
	r.AddCookie(&http.Cookie{Name: "another", Value: "one"})

	ctx := getContext(r)
	assert.Equal(r.Context(), ctx.Context)
	assert.Equal(r.Cookies(), ctx.Cookies)
	assert.Equal("", ctx.XSRFToken)
}

func TestGetContextMissingXSRFToken(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "another", Value: "one"})

	ctx := getContext(r)
	assert.Equal(r.Context(), ctx.Context)
	assert.Equal(r.Cookies(), ctx.Cookies)
	assert.Equal("", ctx.XSRFToken)
}

func TestGetContextForPostRequest(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("POST", "/", strings.NewReader("xsrfToken=the-real-one"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: "z3tVRZ00yx4dHz3KWYv3boLWHZ4/RsCsVAKbvo2SBNc%3D"})
	r.AddCookie(&http.Cookie{Name: "another", Value: "one"})

	ctx := getContext(r)
	assert.Equal(r.Context(), ctx.Context)
	assert.Equal(r.Cookies(), ctx.Cookies)
	assert.Equal("the-real-one", ctx.XSRFToken)
}

func TestCancelledContext(t *testing.T) {
	assert := assert.New(t)

	ctx, logBuf := contextWithLogger()

	client := &mockErrorHandlerClient{}
	tmplError := &mockTemplate{}

	wrap := errorHandler(client, tmplError, "/prefix", "http://sirius")
	handler := wrap(func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		return context.Canceled
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(ctx, "GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(499, resp.StatusCode)
	assert.Equal("", logBuf.String())
}
