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

func TestAddTeam(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		scenario      string
		setup         func()
		name          string
		teamType      string
		phone         string
		email         string
		expectedID    int
		expectedError error
	}{
		{
			scenario: "Created",
			setup: func() {
				pact.
					AddInteraction().
					Given("An admin user").
					UponReceiving("A request to add a new team").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/v1/teams"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"email":       "john.doe@example.com",
							"name":        "testteam",
							"phoneNumber": "0300456090",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusCreated,
						Body: matchers.Like(map[string]interface{}{
							"id": matchers.Like(123),
						}),
					})
			},
			email:      "john.doe@example.com",
			name:       "testteam",
			phone:      "0300456090",
			teamType:   "",
			expectedID: 123,
		},

		{
			scenario: "CreatedSupervision",
			setup: func() {
				pact.
					AddInteraction().
					Given("An admin user").
					UponReceiving("A request to add a new supervision team").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/v1/teams"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"email":       "john.doe@example.com",
							"name":        "supervisiontestteam",
							"phoneNumber": "0300456090",
							"type":        "INVESTIGATIONS",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusCreated,
						Body: matchers.Like(map[string]interface{}{
							"id": matchers.Like(123),
						}),
					})
			},
			email:      "john.doe@example.com",
			name:       "supervisiontestteam",
			phone:      "0300456090",
			teamType:   "INVESTIGATIONS",
			expectedID: 123,
		},
		{
			scenario: "Errors",
			setup: func() {
				pact.
					AddInteraction().
					Given("An admin user").
					UponReceiving("A request to add a new team errors").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/v1/teams"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"email":       "john.doehrfgjuerhujghejrhrgherjrghgjrehergeghrjkrghkerhgerjkhgerjkheghergkhgekrhgerherhjghkjerhgherghjkerhgekjherkjhgerhgjehherkjhgkjehrghrehgkjrehjkghrjkehgrehehgkjhrejghhehgkjerhegjrhegrjhrjkhgkrhrghrkjegrkjehrghjkerhgjkhergjhrjkerregjhrekjhrgrehjkg@example.com",
							"name":        "testteam",
							"phoneNumber": "0300456090",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusBadRequest,
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/problem+json"),
						},
						Body: matchers.Like(map[string]interface{}{
							"validation_errors": matchers.Like(map[string]interface{}{
								"email": matchers.Like(map[string]interface{}{
									"stringLengthTooLong": "The input is more than 255 characters long",
								}),
							}),
						}),
					})
			},
			email:    "john.doehrfgjuerhujghejrhrgherjrghgjrehergeghrjkrghkerhgerjkhgerjkheghergkhgekrhgerherhjghkjerhgherghjkerhgekjherkjhgerhgjehherkjhgkjehrghrehgkjrehjkghrjkehgrehehgkjhrejghhehgkjerhegjrhegrjhrjkhgkrhrghrkjegrkjehrghjkerhgjkhergjhrjkerregjhrekjhrgrehjkg@example.com",
			name:     "testteam",
			phone:    "0300456090",
			teamType: "",
			expectedError: ValidationError{
				Errors: ValidationErrors{
					"email": {
						"stringLengthTooLong": "The input is more than 255 characters long",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				id, err := client.AddTeam(Context{Context: context.Background()}, tc.name, tc.teamType, tc.phone, tc.email)
				assert.Equal(t, tc.expectedError, err)
				assert.Equal(t, tc.expectedID, id)
				return nil
			}))
		})
	}
}

func TestAddTeamStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.AddTeam(Context{Context: context.Background()}, "", "", "", "")
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/v1/teams",
		Method: http.MethodPost,
	}, err)
}
