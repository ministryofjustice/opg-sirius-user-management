package sirius

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser(t *testing.T) {
	pact := &dsl.Pact{
		Consumer:          "sirius-user-management",
		Provider:          "sirius",
		Host:              "localhost",
		PactFileWriteMode: "merge",
		LogDir:            "../../logs",
		PactDir:           "../../pacts",
	}
	defer pact.Teardown()

	testCases := []struct {
		name          string
		setup         func()
		userID        int
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
					WithRequest(dsl.Request{
						Method: http.MethodDelete,
						Path:   dsl.String("/auth/user/123"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

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
			http.Error(w, `{"message":"oops"}`, http.StatusBadRequest)
		}),
	)
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.DeleteUser(Context{Context: context.Background()}, 123)
	assert.Equal(t, ClientError("oops"), err)
}

func TestDeleteUserStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.DeleteUser(Context{Context: context.Background()}, 123)
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/auth/user/123",
		Method: http.MethodDelete,
	}, err)
}
