package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

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

	testCases := []struct {
		name              string
		setup             func()
		expectedMyDetails MyDetails
		expectedError     error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get my details").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/users/current"),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"id":          dsl.Like(47),
							"name":        dsl.Like("system"),
							"phoneNumber": dsl.Like("03004560300"),
							"teams": dsl.EachLike(map[string]interface{}{
								"displayName": dsl.Like("Allocations - (Supervision)"),
							}, 1),
							"displayName": dsl.Like("system admin"),
							"deleted":     dsl.Like(false),
							"email":       dsl.Like("system.admin@opgtest.com"),
							"firstname":   dsl.Like("system"),
							"surname":     dsl.Like("admin"),
							"roles":       dsl.EachLike("System Admin", 1),
							"locked":      dsl.Like(false),
							"suspended":   dsl.Like(false),
						}),
					})
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
				Roles:       []string{"System Admin"},
				Locked:      false,
				Suspended:   false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				myDetails, err := client.MyDetails(Context{Context: context.Background()})
				assert.Equal(t, tc.expectedMyDetails, myDetails)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestMyDetailsStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.MyDetails(Context{Context: context.Background()})
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/users/current",
		Method: http.MethodGet,
	}, err)
}
