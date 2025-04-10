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

func TestEditTeam(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		team          Team
		expectedError func(int) error
	}{
		{
			name: "OK",
			team: Team{
				ID:          65,
				DisplayName: "Test team",
				Type:        "INVESTIGATIONS",
				PhoneNumber: "014729583920",
				Email:       "test.team@opgtest.com",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("Supervision team with members exists").
					UponReceiving("A request to edit the team").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/api/v1/teams/65"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"email":       "test.team@opgtest.com",
							"name":        "Test team",
							"phoneNumber": "014729583920",
							"type":        "INVESTIGATIONS",
							"memberIds":   []int{},
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
					})
			},
			expectedError: func(port int) error { return nil },
		},

		{
			name: "OKSendsMembers",
			team: Team{
				ID:          65,
				DisplayName: "Test team with members",
				Type:        "INVESTIGATIONS",
				PhoneNumber: "014729583920",
				Email:       "test.team@opgtest.com",
				Members: []TeamMember{
					{
						ID:    23,
						Email: "someone@opgtest.com",
					},
				},
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("Supervision team with members exists").
					UponReceiving("A request to edit the team with members").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/api/v1/teams/65"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"email":       "test.team@opgtest.com",
							"name":        "Test team with members",
							"phoneNumber": "014729583920",
							"type":        "INVESTIGATIONS",
							"memberIds":   []int{23},
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
					})
			},
			expectedError: func(port int) error { return nil },
		},
		{
			name: "Validation Errors",
			team: Team{
				ID:          65,
				DisplayName: "Test duplicate finance team",
				Type:        "FINANCE",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("Supervision team with members exists").
					UponReceiving("A request to edit the team with a non-unique type").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/api/v1/teams/65"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"name":        "Test duplicate finance team",
							"type":        "FINANCE",
							"email":       "",
							"phoneNumber": "",
							"memberIds":   []int{},
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusBadRequest,
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/problem+json"),
						},
						Body: matchers.Like(map[string]interface{}{
							"validation_errors": matchers.Like(map[string]interface{}{
								"type": matchers.Like(map[string]interface{}{
									"error": matchers.Like("Invalid team type"),
								}),
							}),
						}),
					})
			},
			expectedError: func(port int) error {
				return &ValidationError{
					Errors: ValidationErrors{
						"type": {
							"error": "Invalid team type",
						},
					},
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.EditTeam(Context{Context: context.Background()}, tc.team)

				assert.Equal(t, tc.expectedError(config.Port), err)
				return nil
			}))
		})
	}
}

func TestEditTeamStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.EditTeam(Context{Context: context.Background()}, Team{ID: 65})
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/teams/65",
		Method: http.MethodPut,
	}, err)
}
