package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestAddTeam(t *testing.T) {
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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/api/v1/teams"),
						Body: map[string]interface{}{
							"email":       "john.doe@example.com",
							"name":        "testteam",
							"phoneNumber": "0300456090",
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusCreated,
						Body: dsl.Like(map[string]interface{}{
							"id": dsl.Like(123),
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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/api/v1/teams"),
						Body: map[string]interface{}{
							"email":       "john.doe@example.com",
							"name":        "supervisiontestteam",
							"phoneNumber": "0300456090",
							"type":        "INVESTIGATIONS",
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusCreated,
						Body: dsl.Like(map[string]interface{}{
							"id": dsl.Like(123),
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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/api/v1/teams"),
						Body: map[string]interface{}{
							"email":       "john.doehrfgjuerhujghejrhrgherjrghgjrehergeghrjkrghkerhgerjkhgerjkheghergkhgekrhgerherhjghkjerhgherghjkerhgekjherkjhgerhgjehherkjhgkjehrghrehgkjrehjkghrjkehgrehehgkjhrejghhehgkjerhegjrhegrjhrjkhgkrhrghrkjegrkjehrghjkerhgjkhergjhrjkerregjhrekjhrgrehjkg@example.com",
							"name":        "testteam",
							"phoneNumber": "0300456090",
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusBadRequest,
						Body: dsl.Like(map[string]interface{}{
							"validation_errors": dsl.Like(map[string]interface{}{
								"email": dsl.Like(map[string]interface{}{
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

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

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
		URL:    s.URL + "/api/v1/teams",
		Method: http.MethodPost,
	}, err)
}
