package sirius

import (
	"context"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-user-management/internal/model"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestAddFeedback(t *testing.T) {
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
		email         string
		caseNumber    string
		feedback      string
		expectedError error
	}{
		{
			scenario: "Created",
			setup: func() {
				pact.
					AddInteraction().
					Given("An user").
					UponReceiving("A request to add feedback").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/api/supervision-feedback"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"name":       "testteam",
							"email":      "john.doe@example.com",
							"caseNumber": "12345",
							"feedback":   "something looks bad",
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusCreated,
					})
			},
			name:       "testteam",
			email:      "john.doe@example.com",
			caseNumber: "12345",
			feedback:   "something looks bad",
		},
		{
			scenario: "Errors",
			setup: func() {
				pact.
					AddInteraction().
					Given("An user").
					UponReceiving("A request to add feedback errors").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/api/supervision-feedback"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"name":       "testteam",
							"email":      "john.doe@example.com",
							"caseNumber": "12345",
							"feedback":   "rfgjuerhujghejrhrgherjrghgjrehergeghrjkrghkerhgerjkhgerjkheghergkhgekrhgerherhjghkjerhgherghjkerhgekjherkjhgerhgjehherkjhgkjehrghrehgkjrehjkghrjkehgrehehgkjhrejghhehgkjerhegjrhegrjhrjkhgkrhrghrkjegrkjehrghjkerhgjkhergjhrjkerregjhrekjhrgrehjkgdoehrfgjuerhujghejrhrgherjrghgjrehergeghrjkrghkerhgerjkhgerjkheghergkhgekrhgerherhjghkjerhgherghjkerhgekjherkjhgerhgjehherkjhgkjehrghrehgkjrehjkghrjkehgrehehgkjhrejghhehgkjerhegjrhegrjhrjkhgkrhrghrkjegrkjehrghjkerhgjkhergjhrjkerregjhrekjhrgrehjkgrfgjuerhujghejrhrgherjrghgjrehergeghrjkrghkerhgerjkhgerjkheghergkhgekrhgerherhjghkjerhgherghjkerhgekjherkjhgerhgjehherkjhgkjehrghrehgkjrehjkghrjkehgrehehgkjhrejghhehgkjerhegjrhegrjhrjkhgkrhrghrkjegrkjehrghjkerhgjkhergjhrjkerregjhrekjhrgrehjkgrfgjuerhujghejrhrgherjrghgjrehergeghrjkrghkerhgerjkhgerjkheghergkhgekrhgerherhjghkjerhgherghjkerhgekjherkjhgerhgjehherkjhgkjehrghrehgkjrehjkghrjkehgrehehgkjhrejghhehgkjerhegjrhegrjhrjkhgkrhrghrkjegrkjehrghjkerhgjkhergjhrjkerregjhrekjhrgrehjkgrfgjuerhujghejrhrgherjrghgjrehergeghrjkrghkerhgerjkhgerjkheghergkhgekrhgerherhjghkjerhgherghjkerhgekjherkjhgerhgjehherkjhgkjehrghrehgkjrehjkghrjkehgrehehgkjhrejghhehgkjerhegjrhegrjhrjkhgkrhrghrkjegrkjehrghjkerhgjkhergjhrjkerregjhrekjhrgrehjkg",
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusCreated,
						Body: dsl.Like(map[string]interface{}{
							"validation_errors": dsl.Like(map[string]interface{}{
								"email": dsl.Like(map[string]interface{}{
									"stringLengthTooLong": "The input is more than 900 characters long",
								}),
							}),
						}),
					})
			},
			name:       "testteam",
			email:      "john.doe@example.com",
			caseNumber: "12345",
			feedback:   "rfgjuerhujghejrhrgherjrghgjrehergeghrjkrghkerhgerjkhgerjkheghergkhgekrhgerherhjghkjerhgherghjkerhgekjherkjhgerhgjehherkjhgkjehrghrehgkjrehjkghrjkehgrehehgkjhrejghhehgkjerhegjrhegrjhrjkhgkrhrghrkjegrkjehrghjkerhgjkhergjhrjkerregjhrekjhrgrehjkgdoehrfgjuerhujghejrhrgherjrghgjrehergeghrjkrghkerhgerjkhgerjkheghergkhgekrhgerherhjghkjerhgherghjkerhgekjherkjhgerhgjehherkjhgkjehrghrehgkjrehjkghrjkehgrehehgkjhrejghhehgkjerhegjrhegrjhrjkhgkrhrghrkjegrkjehrghjkerhgjkhergjhrjkerregjhrekjhrgrehjkgrfgjuerhujghejrhrgherjrghgjrehergeghrjkrghkerhgerjkhgerjkheghergkhgekrhgerherhjghkjerhgherghjkerhgekjherkjhgerhgjehherkjhgkjehrghrehgkjrehjkghrjkehgrehehgkjhrejghhehgkjerhegjrhegrjhrjkhgkrhrghrkjegrkjehrghjkerhgjkhergjhrjkerregjhrekjhrgrehjkgrfgjuerhujghejrhrgherjrghgjrehergeghrjkrghkerhgerjkhgerjkheghergkhgekrhgerherhjghkjerhgherghjkerhgekjherkjhgerhgjehherkjhgkjehrghrehgkjrehjkghrjkehgrehehgkjhrejghhehgkjerhegjrhegrjhrjkhgkrhrghrkjegrkjehrghjkerhgjkhergjhrjkerregjhrekjhrgrehjkgrfgjuerhujghejrhrgherjrghgjrehergeghrjkrghkerhgerjkhgerjkheghergkhgekrhgerherhjghkjerhgherghjkerhgekjherkjhgerhgjehherkjhgkjehrghrehgkjrehjkghrjkehgrehehgkjhrejghhehgkjerhegjrhegrjhrjkhgkrhrghrkjegrkjehrghjkerhgjkhergjhrjkerregjhrekjhrgrehjkg",
			expectedError: ValidationError{
				Errors: ValidationErrors{
					"email": {
						"stringLengthTooLong": "The input is more than 900 characters long",
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

				err := client.AddFeedback(Context{Context: context.Background()}, model.FeedbackForm{
					IsSupervisionFeedback: false,
					Name:                  tc.name,
					Email:                 tc.email,
					CaseNumber:            tc.caseNumber,
					Message:               tc.feedback,
				})
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestAddFeedbackIsEmptyValidationError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.AddFeedback(Context{Context: context.Background()}, model.FeedbackForm{
		IsSupervisionFeedback: true,
		Name:                  "Toad",
		Email:                 "toad@toadhall.com",
		CaseNumber:            "123",
		Message:               "",
	})
	assert.Equal(t, ValidationError{
		Message: "isEmpty",
		Errors:  nil,
	}, err)
}

func TestAddFeedbackStringTooLongValidationError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.AddFeedback(Context{Context: context.Background()}, model.FeedbackForm{
		IsSupervisionFeedback: true,
		Name:                  "Toad",
		Email:                 "toad@toadhall.com",
		CaseNumber:            "123",
		Message:               strings.Repeat("a", 901),
	})
	assert.Equal(t, ValidationError{
		Message: "stringLengthTooLong",
		Errors:  nil,
	}, err)
}
