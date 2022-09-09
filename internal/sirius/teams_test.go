package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestTeamsIgnoredPact(t *testing.T) {
	// These tests are pretty much impossible to validate in Sirius due to the
	// responses containing a mix of LPA and Supervision teams, and Pact not
	// having a way to match arrays containing either type. We still want to
	// produce a contract that can be used for mocking responses though.
	pact := &dsl.Pact{
		Consumer:          "ignored",
		Provider:          "ignored",
		Host:              "localhost",
		PactFileWriteMode: "merge",
		LogDir:            "../../logs",
		PactDir:           "../../pacts",
	}
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams"),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.EachLike(map[string]interface{}{
							"id":          dsl.Like(123),
							"displayName": dsl.Like("Cool Team"),
							"members": dsl.EachLike(map[string]interface{}{
								"displayName": dsl.Like("John"),
								"email":       dsl.Like("john@opgtest.com"),
							}, 1),
							"teamType": dsl.Like(map[string]interface{}{
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams"),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.EachLike(map[string]interface{}{
							"id":          dsl.Like(123),
							"displayName": dsl.Like("Cool Team"),
							"members": dsl.EachLike(map[string]interface{}{
								"displayName": dsl.Like("John"),
								"email":       dsl.Like("john@opgtest.com"),
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

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

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
