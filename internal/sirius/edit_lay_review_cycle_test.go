package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

type editLayReviewCycleBadRequestResponse struct {
	Status           int    `json:"status" pact:"example=400"`
	Detail           string `json:"detail" pact:"example=Enter a review cycle between 1 and 10 for lay cases"`
}

func TestEditLayReviewCycle(t *testing.T) {
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
		name                    string
		layPercentage			string
		reviewCycle				string
		setup                   func()
		cookies                 []*http.Cookie
		expectedError 			error
	}{
		{
			name: "Validation errors",
			reviewCycle: "15",
			layPercentage: "50",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to edit the lay review cycle errors on validation").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/api/v1/random-review-settings"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
						Body: map[string]interface{}{
                            "reviewCycle": "15",
							"layPercentage": "50",
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusBadRequest,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/problem+json")},
						Body: dsl.Match(editLayReviewCycleBadRequestResponse{}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedError: ValidationError{
				Message: "Enter a review cycle between 1 and 10 for lay cases",
			},
		},
		{
			name: "OK",
			reviewCycle: "1",
			layPercentage: "30",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to edit the lay review cycle").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/api/v1/random-review-settings"),
						Headers: dsl.MapMatcher{
							"Content-type":        dsl.String("application/json"),
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
						Body: map[string]interface{}{
							"reviewCycle": "1",
							"layPercentage": "30",
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedError: nil,
		},
		{
			name: "Unauthorized",
			reviewCycle: "9",
			layPercentage: "15",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get lay review cycle without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/api/v1/random-review-settings"),
						Headers: dsl.MapMatcher{
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
						Body: map[string]interface{}{
							"reviewCycle": "9",
							"layPercentage": "15",
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: ErrUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.EditLayReviewCycle(getContext(tc.cookies), tc.reviewCycle, tc.layPercentage)

				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestEditLayReviewCycleStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.EditLayReviewCycle(getContext(nil), "15", "9")
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/random-review-settings",
		Method: http.MethodPost,
	}, err)
}
