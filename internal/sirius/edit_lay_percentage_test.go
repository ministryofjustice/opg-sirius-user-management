package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

type editLayPercentageBadRequestResponse struct {
	Status           int    `json:"status" pact:"example=400"`
	Detail           string `json:"detail" pact:"example=Enter a percentage between 0 and 100 for lay cases"`
}

const UserExists = "User exists";
const UrlRoute = "/api/v1/random-review-settings";
const BypassMembrane = "OPG-Bypass-Membrane";

func TestEditLayPercentage(t *testing.T) {
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
			layPercentage: "200",
			reviewCycle: "3",
			setup: func() {
				pact.
					AddInteraction().
					Given(UserExists).
					UponReceiving("A request to edit the lay percentage errors on validation").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String(UrlRoute),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							BypassMembrane: dsl.String("1"),
						},
						Body: map[string]interface{}{
							"layPercentage": "200",
							"reviewCycle": "3",
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusBadRequest,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/problem+json")},
						Body: dsl.Match(editLayPercentageBadRequestResponse{}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedError: ValidationError{
				Message: "Enter a percentage between 0 and 100 for lay cases",
			},
		},
		{
			name: "OK",
			layPercentage: "20",
			reviewCycle: "3",
			setup: func() {
				pact.
					AddInteraction().
					Given(UserExists).
					UponReceiving("A request to edit the lay percentage").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String(UrlRoute),
						Headers: dsl.MapMatcher{
							"Content-type":        dsl.String("application/json"),
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							BypassMembrane: dsl.String("1"),
						},
						Body: map[string]interface{}{
							"layPercentage": "20",
							"reviewCycle": "3",
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
			layPercentage: "20",
			reviewCycle: "3",
			setup: func() {
				pact.
					AddInteraction().
					Given(UserExists).
					UponReceiving("A request to get lay percentage without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String(UrlRoute),
						Headers: dsl.MapMatcher{
							BypassMembrane: dsl.String("1"),
						},
						Body: map[string]interface{}{
							"layPercentage": "20",
							"reviewCycle": "3",
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

				err := client.EditLayPercentageReviewCycle(getContext(tc.cookies), tc.reviewCycle, tc.layPercentage)

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

	err := client.EditLayPercentageReviewCycle(getContext(nil), "3", "20")
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + UrlRoute,
		Method: http.MethodPost,
	}, err)
}
