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

const UserExists = "User exists"
const UrlRoute = "/supervision-api/v1/random-review-settings"

func TestEditLayPercentage(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		layPercentage string
		paPercentage  string
		proPercentage string
		reviewCycle   string
		setup         func()
		expectedError error
	}{
		{
			name:          "Lay percentage - validation errors",
			layPercentage: "200",
			paPercentage:  "30",
			proPercentage: "0",
			reviewCycle:   "3",
			setup: func() {
				pact.
					AddInteraction().
					Given(UserExists).
					UponReceiving("A request to edit the lay percentage errors on validation").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String(UrlRoute),
						Body: map[string]interface{}{
							"layPercentage": "200",
							"paPercentage":  "30",
							"proPercentage": "0",
							"reviewCycle":   "3",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusBadRequest,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/problem+json")},
						Body: matchers.Like(map[string]interface{}{
							"status": matchers.Like(400),
							"detail": matchers.Like("Enter a percentage between 0 and 100 for lay cases"),
						}),
					})
			},
			expectedError: ValidationError{
				Message: "Enter a percentage between 0 and 100 for lay cases",
			},
		},
		{
			name:          "Lay percentage - success",
			layPercentage: "20",
			paPercentage:  "10",
			proPercentage: "18",
			reviewCycle:   "3",
			setup: func() {
				pact.
					AddInteraction().
					Given(UserExists).
					UponReceiving("A request to edit the lay percentage").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String(UrlRoute),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"layPercentage": "20",
							"paPercentage":  "10",
							"proPercentage": "18",
							"reviewCycle":   "3",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
					})
			},
			expectedError: nil,
		},
		{
			name:          "PA percentage - validation errors",
			layPercentage: "20",
			paPercentage:  "1000",
			proPercentage: "18",
			reviewCycle:   "3",
			setup: func() {
				pact.
					AddInteraction().
					Given(UserExists).
					UponReceiving("A request to edit the PA percentage errors on validation").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String(UrlRoute),
						Body: map[string]interface{}{
							"layPercentage": "20",
							"paPercentage":  "1000",
							"proPercentage": "18",
							"reviewCycle":   "3",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusBadRequest,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/problem+json")},
						Body: matchers.Like(map[string]interface{}{
							"status": matchers.Like(400),
							"detail": matchers.Like("Enter a percentage between 0 and 100 for PA cases"),
						}),
					})
			},
			expectedError: ValidationError{
				Message: "Enter a percentage between 0 and 100 for PA cases",
			},
		},
		{
			name:          "PA percentage - success",
			layPercentage: "20",
			paPercentage:  "50",
			proPercentage: "18",
			reviewCycle:   "3",
			setup: func() {
				pact.
					AddInteraction().
					Given(UserExists).
					UponReceiving("A request to edit the PA percentage").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String(UrlRoute),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"layPercentage": "20",
							"paPercentage":  "50",
							"proPercentage": "18",
							"reviewCycle":   "3",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
					})
			},
			expectedError: nil,
		},
		{
			name:          "PRO percentage - validation errors",
			layPercentage: "20",
			paPercentage:  "50",
			proPercentage: "2000",
			reviewCycle:   "3",
			setup: func() {
				pact.
					AddInteraction().
					Given(UserExists).
					UponReceiving("A request to edit the PRO percentage errors on validation").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String(UrlRoute),
						Body: map[string]interface{}{
							"layPercentage": "20",
							"paPercentage":  "50",
							"proPercentage": "2000",
							"reviewCycle":   "3",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusBadRequest,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/problem+json")},
						Body: matchers.Like(map[string]interface{}{
							"status": matchers.Like(400),
							"detail": matchers.Like("Enter a percentage between 0 and 100 for Pro cases"),
						}),
					})
			},
			expectedError: ValidationError{
				Message: "Enter a percentage between 0 and 100 for Pro cases",
			},
		},
		{
			name:          "PRO percentage - success",
			layPercentage: "20",
			paPercentage:  "50",
			proPercentage: "18",
			reviewCycle:   "3",
			setup: func() {
				pact.
					AddInteraction().
					Given(UserExists).
					UponReceiving("A request to edit the PRO percentage").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String(UrlRoute),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"layPercentage": "20",
							"paPercentage":  "50",
							"proPercentage": "18",
							"reviewCycle":   "3",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
					})
			},
			expectedError: nil,
		},
		{
			name:          "Review cycle - success",
			layPercentage: "20",
			paPercentage:  "50",
			proPercentage: "18",
			reviewCycle:   "5",
			setup: func() {
				pact.
					AddInteraction().
					Given(UserExists).
					UponReceiving("A request to edit the review cycle").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String(UrlRoute),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"layPercentage": "20",
							"paPercentage":  "50",
							"proPercentage": "18",
							"reviewCycle":   "5",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
					})
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))
				data := EditRandomReview{tc.layPercentage, tc.paPercentage, tc.proPercentage, tc.reviewCycle}

				err := client.EditRandomReviewSettings(Context{Context: context.Background()}, data)

				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestEditLayPercentageStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.EditRandomReviewSettings(Context{Context: context.Background()}, EditRandomReview{"3", "10", "18", "20"})
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + UrlRoute,
		Method: http.MethodPost,
	}, err)
}
