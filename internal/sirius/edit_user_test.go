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
	Message string `json:"message" pact:"example=oops"`
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
		cookies       []*http.Cookie
		expectedError func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("A user").
					UponReceiving("A request to edit the user").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/auth/user/123"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
						Body: map[string]interface{}{
							"id":        123,
							"firstname": "a",
							"surname":   "b",
							"roles":     []string{"e", "f", "d"},
							"locked":    false,
							"suspended": true,
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedError: func(port int) error { return nil },
		},

		{
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("A user").
					UponReceiving("A request edit the user without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/auth/user/123"),
						Headers: dsl.MapMatcher{
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: func(port int) error { return ErrUnauthorized },
		},

		{
			name: "Validation Errors",
			setup: func() {
				pact.
					AddInteraction().
					Given("A user").
					UponReceiving("A request to edit the user errors on validation").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/auth/user/123"),
						Headers: dsl.MapMatcher{
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusBadRequest,
						Body:   dsl.Match(editUserErrorsResponse{}),
					})
			},
			expectedError: func(port int) error { return ClientError("oops") },
		},

		{
			name: "Errors",
			setup: func() {
				pact.
					AddInteraction().
					Given("A user").
					UponReceiving("A request to edit the user errors").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/auth/user/123"),
						Headers: dsl.MapMatcher{
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusBadRequest,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:   http.StatusBadRequest,
					URL:    fmt.Sprintf("http://localhost:%d/auth/user/123", port),
					Method: http.MethodPut,
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.EditUser(context.Background(), tc.cookies, AuthUser{
					ID:           123,
					Firstname:    "a",
					Surname:      "b",
					Organisation: "d",
					Roles:        []string{"e", "f"},
					Locked:       false,
					Suspended:    true,
				})

				assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				return nil
			}))
		})
	}
}
