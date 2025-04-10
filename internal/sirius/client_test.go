package sirius

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/stretchr/testify/assert"
)

func teapotServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		}),
	)
}

func invalidJSONServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("1a is not valid json"))
		}),
	)
}

func newPact() (*consumer.V4HTTPMockProvider, error) {
	return consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-user-management",
		Provider: "sirius",
		Host:     "127.0.0.1",
		LogDir:   "../../logs",
		PactDir:  "../../pacts",
	})
}

func newIgnoredPact() (*consumer.V4HTTPMockProvider, error) {
	return consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "ignored",
		Provider: "ignored",
		Host:     "127.0.0.1",
		LogDir:   "../../logs",
		PactDir:  "../../pacts",
	})
}

func TestClientError(t *testing.T) {
	assert.Equal(t, "message", ClientError("message").Error())
}

func TestValidationError(t *testing.T) {
	assert.Equal(t, "message", ValidationError{Message: "message"}.Error())
}

func TestStatusError(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/some/url", nil)

	resp := &http.Response{
		StatusCode: http.StatusTeapot,
		Request:    req,
	}

	err := newStatusError(resp)

	assert.Equal(t, "POST /some/url returned 418", err.Error())
	assert.Equal(t, "unexpected response from Sirius", err.Title())
	assert.Equal(t, err, err.Data())
}
