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

func TestTeam(t *testing.T) {
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
		id               int
		name             string
		setup            func()
		expectedResponse Team
		expectedError    error
	}{
		{
			name: "OK",
			id:   65,
			setup: func() {
				pact.
					AddInteraction().
					Given("Supervision team with members exists").
					UponReceiving("A request for a team").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams/65"),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"id":          dsl.Like(65),
							"displayName": dsl.Like("Cool Team"),
							"email":       dsl.Like("coolteam@opgtest.com"),
							"phoneNumber": dsl.Like("01818118181"),
							"members": dsl.EachLike(map[string]interface{}{
								"displayName": dsl.Like("John"),
								"email":       dsl.Like("john@opgtest.com"),
							}, 1),
							"teamType": dsl.Like(map[string]interface{}{
								"handle": "ALLOCATIONS",
							}),
						}),
					})
			},
			expectedResponse: Team{
				ID:          65,
				DisplayName: "Cool Team",
				Email:       "coolteam@opgtest.com",
				PhoneNumber: "01818118181",
				Members: []TeamMember{
					{
						DisplayName: "John",
						Email:       "john@opgtest.com",
					},
				},
				Type: "ALLOCATIONS",
			},
		},
		{
			name: "OKWithLpaTeams",
			id:   65,
			setup: func() {
				pact.
					AddInteraction().
					Given("LPA team with members exists").
					UponReceiving("A request for an LPA team").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams/65"),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"id":          dsl.Like(65),
							"displayName": dsl.Like("Cool Team"),
							"members": dsl.EachLike(map[string]interface{}{
								"displayName": dsl.Like("Carline"),
								"email":       dsl.Like("carline@opgtest.com"),
							}, 1),
						}),
					})
			},
			expectedResponse: Team{
				ID:          65,
				DisplayName: "Cool Team",
				Members: []TeamMember{
					{
						DisplayName: "Carline",
						Email:       "carline@opgtest.com",
					},
				},
				Type: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				team, err := client.Team(Context{Context: context.Background()}, tc.id)
				assert.Equal(t, tc.expectedResponse, team)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestTeamBadJSONResponse(t *testing.T) {
	s := invalidJSONServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.Team(Context{Context: context.Background()}, 123)
	assert.IsType(t, &json.UnmarshalTypeError{}, err)
}

func TestTeamStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.Team(Context{Context: context.Background()}, 123)
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/teams/123",
		Method: http.MethodGet,
	}, err)
}
