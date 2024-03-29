package sirius

import (
	"context"
	"encoding/json"
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

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse AuthUser
		expectedError    error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request for the user").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/users/123"),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body:    dsl.Match(&exampleAuthUser{}),
					})
			},
			expectedResponse: AuthUser{
				ID:           123,
				Firstname:    "system",
				Surname:      "admin",
				Email:        "system.admin@opgtest.com",
				Organisation: "",
				Roles:        []string{"string"},
				Suspended:    false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				users, err := client.User(Context{Context: context.Background()}, 123)
				assert.Equal(t, tc.expectedResponse, users)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestUserBadJSONResponse(t *testing.T) {
	s := invalidJSONServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.User(Context{Context: context.Background()}, 123)
	assert.IsType(t, &json.UnmarshalTypeError{}, err)
}

func TestUserStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.User(Context{Context: context.Background()}, 123)
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/users/123",
		Method: http.MethodGet,
	}, err)
}
