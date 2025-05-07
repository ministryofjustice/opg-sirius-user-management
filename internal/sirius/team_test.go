package sirius

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestTeam(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/v1/teams/65"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"id":          matchers.Like(65),
							"displayName": matchers.Like("Cool Team"),
							"email":       matchers.Like("coolteam@opgtest.com"),
							"phoneNumber": matchers.Like("01818118181"),
							"members": matchers.EachLike(map[string]interface{}{
								"displayName": matchers.Like("John"),
								"email":       matchers.Like("john@opgtest.com"),
							}, 1),
							"teamType": matchers.Like(map[string]interface{}{
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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/v1/teams/65"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"id":          matchers.Like(65),
							"displayName": matchers.Like("Cool Team"),
							"members": matchers.EachLike(map[string]interface{}{
								"displayName": matchers.Like("Carline"),
								"email":       matchers.Like("carline@opgtest.com"),
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

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

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
		URL:    s.URL + "/v1/teams/123",
		Method: http.MethodGet,
	}, err)
}
