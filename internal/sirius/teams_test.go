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

func TestTeamsIgnoredPact(t *testing.T) {
	// These tests are pretty much impossible to validate in Sirius due to the
	// responses containing a mix of LPA and Supervision teams, and Pact not
	// having a way to match arrays containing either type. We still want to
	// produce a contract that can be used for mocking responses though.
	pact, err := newIgnoredPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []Team
		expectedError    error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request for teams").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/teams"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.EachLike(map[string]interface{}{
							"id":          matchers.Like(123),
							"displayName": matchers.Like("Cool Team"),
							"members": matchers.EachLike(map[string]interface{}{
								"displayName": matchers.Like("John"),
								"email":       matchers.Like("john@opgtest.com"),
							}, 1),
							"teamType": matchers.Like(map[string]interface{}{
								"handle": "ALLOCATIONS",
								"label":  "Allocations",
							}),
						}, 1),
					})
			},
			expectedResponse: []Team{
				{
					ID:          123,
					DisplayName: "Cool Team",
					Members: []TeamMember{
						{
							DisplayName: "John",
							Email:       "john@opgtest.com",
						},
					},
					Type:      "ALLOCATIONS",
					TypeLabel: "Supervision — Allocations",
				},
			},
		},

		{
			name: "OKWithLpaTeams",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists and teams have no type").
					UponReceiving("A request for teams").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/teams"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.EachLike(map[string]interface{}{
							"id":          matchers.Like(123),
							"displayName": matchers.Like("Cool Team"),
							"members": matchers.EachLike(map[string]interface{}{
								"displayName": matchers.Like("John"),
								"email":       matchers.Like("john@opgtest.com"),
							}, 1),
						}, 1),
					})
			},
			expectedResponse: []Team{
				{
					ID:          123,
					DisplayName: "Cool Team",
					Members: []TeamMember{
						{
							DisplayName: "John",
							Email:       "john@opgtest.com",
						},
					},
					Type:      "",
					TypeLabel: "LPA",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				users, err := client.Teams(Context{Context: context.Background()})
				assert.Equal(t, tc.expectedResponse, users)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestTeamsStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.Teams(Context{Context: context.Background()})
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/teams",
		Method: http.MethodGet,
	}, err)
}
