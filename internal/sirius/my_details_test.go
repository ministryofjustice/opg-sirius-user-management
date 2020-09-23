package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

type myDetailsResponse struct {
	ID          int    `json:"id" pact:"example=47"`
	Name        string `json:"name" pact:"example=system"`
	PhoneNumber string `json:"phoneNumber" pact:"example=03004560300"`
	Teams       []struct {
		DisplayName string `json:"displayName" pact:"example=Allocations - (Supervision)"`
	} `json:"teams"`
	DisplayName string   `json:"displayName" pact:"example=system admin"`
	Deleted     bool     `json:"deleted" pact:"example=false"`
	Email       string   `json:"email" pact:"example=system.admin@opgtest.com"`
	Firstname   string   `json:"firstname" pact:"example=system"`
	Surname     string   `json:"surname" pact:"example=admin"`
	Roles       []string `json:"roles"`
	Locked      bool     `json:"locked" pact:"example=false"`
	Suspended   bool     `json:"suspended" pact:"example=false"`
}

func TestMyDetails(t *testing.T) {
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
		setup             func()
		cookies           []*http.Cookie
		expectedMyDetails MyDetails
		expectedError     error
	}{
		"OK": {
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get my details").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/users/current"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body:    dsl.Match(myDetailsResponse{}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedMyDetails: MyDetails{
				ID:          47,
				Name:        "system",
				PhoneNumber: "03004560300",
				Teams: []MyDetailsTeam{
					{DisplayName: "Allocations - (Supervision)"},
				},
				DisplayName: "system admin",
				Deleted:     false,
				Email:       "system.admin@opgtest.com",
				Firstname:   "system",
				Surname:     "admin",
				Roles:       []string{"string"},
				Locked:      false,
				Suspended:   false,
			},
		},

		"Unauthorized": {
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get my details without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/users/current"),
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

				myDetails, err := client.MyDetails(context.Background(), tc.cookies)
				assert.Equal(t, tc.expectedMyDetails, myDetails)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}
