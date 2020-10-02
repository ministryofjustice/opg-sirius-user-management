package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

type exampleAuthUser struct {
	ID        int      `json:"id" pact:"example=123"`
	Firstname string   `json:"firstname" pact:"example=system"`
	Surname   string   `json:"surname" pact:"example=admin"`
	Email     string   `json:"email" pact:"example=system.admin@opgtest.com"`
	Roles     []string `json:"roles"`
	Locked    bool     `json:"locked" pact:"example=true"`
	Suspended bool     `json:"suspended" pact:"example=false"`
}

func TestUser(t *testing.T) {
	pact := &dsl.Pact{
		Consumer:          "sirius-user-management",
		Provider:          "sirius",
		Host:              "localhost",
		PactFileWriteMode: "merge",
		LogDir:            "../../logs",
		PactDir:           "../../pacts",
	}
	defer pact.Teardown()

	testCases := map[string]struct {
		setup            func()
		cookies          []*http.Cookie
		expectedResponse AuthUser
		expectedError    error
	}{
		"OK": {
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request for the user").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/auth/user/123"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body:    dsl.Match(&exampleAuthUser{}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: AuthUser{
				ID:           123,
				Firstname:    "system",
				Surname:      "admin",
				Email:        "system.admin@opgtest.com",
				Organisation: "",
				Roles:        []string{"string"},
				Locked:       true,
				Suspended:    false,
			},
		},

		"Unauthorized": {
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request for the user without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/auth/user/123"),
						Headers: dsl.MapMatcher{
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: ErrUnauthorized,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				users, err := client.User(context.Background(), tc.cookies, 123)
				assert.Equal(t, tc.expectedResponse, users)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}
