package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

type exampleUser struct {
	ID          int    `json:"id" pact:"example=47"`
	DisplayName string `json:"displayName" pact:"example=system admin"`
	Surname     string `json:"surname" pact:"example=admin"`
	Email       string `json:"email" pact:"example=system.admin@opgtest.com"`
	Locked      bool   `json:"locked" pact:"example=false"`
	Suspended   bool   `json:"suspended" pact:"example=false"`
}

func TestListUsers(t *testing.T) {
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
		expectedResponse []User
		expectedError    error
	}{
		"OK": {
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get all users").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/users"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body:    dsl.Match([]exampleUser{}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: []User{
				{
					ID:          47,
					DisplayName: "system admin",
					Email:       "system.admin@opgtest.com",
					Status:      "Active",
				},
			},
		},
		"Unauthorized": {
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get all users without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/users"),
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

				users, err := client.ListUsers(context.Background(), tc.cookies)
				assert.Equal(t, tc.expectedResponse, users)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}
