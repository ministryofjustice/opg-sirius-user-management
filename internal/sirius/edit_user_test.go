package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

type editUserErrorsResponse struct {
	Errors *struct {
		Firstname *struct {
			TooLong string `json:"stringLengthTooLong" pact:"example=First name must be 255 characters or fewer"`
		} `json:"firstname"`
	} `json:"validation_errors"`
}

func TestEditUser(t *testing.T) {
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
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/api/v1/users/123"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
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
					WillRespondWith(dsl.Response{
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
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/api/v1/users/123"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"id":        123,
							"firstname": "grehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjger",
							"surname":   "",
							"roles":     []string{""},
							"suspended": false,
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusBadRequest,
						Body:   dsl.Match(editUserErrorsResponse{}),
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

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

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
		URL:    s.URL + "/api/v1/users/123",
		Method: http.MethodPut,
	}, err)
}
