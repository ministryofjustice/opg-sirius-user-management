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

func TestMyDetails(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String(SupervisionAPIPath + "/v1/users/current"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"id":          matchers.Like(47),
							"name":        matchers.Like("system"),
							"phoneNumber": matchers.Like("03004560300"),
							"teams": matchers.EachLike(map[string]interface{}{
								"displayName": matchers.Like("Allocations - (Supervision)"),
							}, 1),
							"displayName": matchers.Like("system admin"),
							"deleted":     matchers.Like(false),
							"email":       matchers.Like("system.admin@opgtest.com"),
							"firstname":   matchers.Like("system"),
							"surname":     matchers.Like("admin"),
							"roles":       matchers.EachLike("System Admin", 1),
							"suspended":   matchers.Like(false),
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
				Suspended:   false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

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
		URL:    s.URL + SupervisionAPIPath + "/v1/users/current",
		Method: http.MethodGet,
	}, err)
}
