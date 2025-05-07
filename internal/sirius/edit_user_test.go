package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestEditUser(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		user          AuthUser
		expectedError error
	}{
		{
			name: "OK",
			user: AuthUser{
				ID:           123,
				Email:        "c@opgtest.com",
				Firstname:    "a",
				Surname:      "b",
				Organisation: "d",
				Roles:        []string{"e", "f"},
				Suspended:    true,
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to edit the user").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String(SupervisionAPIPath + "/v1/users/123"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"id":        123,
							"email":     "c@opgtest.com",
							"firstname": "a",
							"surname":   "b",
							"roles":     []string{"e", "f", "d"},
							"suspended": true,
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
					})
			},
		},
		{
			name: "Validation Errors",
			user: AuthUser{
				ID:        123,
				Firstname: "grehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjger",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to edit the user errors on validation").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String(SupervisionAPIPath + "/v1/users/123"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"id":        123,
							"firstname": "grehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjger",
							"surname":   "",
							"roles":     []string{""},
							"suspended": false,
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusBadRequest,
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/problem+json"),
						},
						Body: matchers.Like(map[string]interface{}{
							"validation_errors": matchers.Like(map[string]interface{}{
								"firstname": matchers.Like(map[string]interface{}{
									"stringLengthTooLong": matchers.Like("First name must be 255 characters or fewer"),
								}),
							}),
						}),
					})
			},
			expectedError: ValidationError{
				Errors: ValidationErrors{
					"firstname": {
						"stringLengthTooLong": "First name must be 255 characters or fewer",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.EditUser(Context{Context: context.Background()}, tc.user)

				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestEditUserStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.EditUser(Context{Context: context.Background()}, AuthUser{ID: 123})
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + SupervisionAPIPath + "/v1/users/123",
		Method: http.MethodPut,
	}, err)
}
