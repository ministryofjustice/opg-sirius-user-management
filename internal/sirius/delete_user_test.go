package sirius

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		userID        int
		cookies       []*http.Cookie
		expectedError error
	}{
		{
			name:   "OK",
			userID: 123,
			setup: func() {
				pact.
					AddInteraction().
					Given("A user that can be deleted").
					UponReceiving("A request to delete the user").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodDelete,
						Path:   matchers.String("/api/v1/users/123"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.DeleteUser(Context{Context: context.Background()}, tc.userID)

				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestDeleteUserClientError(t *testing.T) {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `{"detail":"oops"}`, http.StatusBadRequest)
		}),
	)
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.DeleteUser(Context{Context: context.Background()}, 123)
	assert.Equal(t, ValidationError{
		Message: "oops",
		Errors: ValidationErrors{
			"#": {"error": "oops"},
		},
	}, err)
}

func TestDeleteUserStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.DeleteUser(Context{Context: context.Background()}, 123)
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/users/123",
		Method: http.MethodDelete,
	}, err)
}
