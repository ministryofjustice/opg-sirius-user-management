package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

type editLayPercentageBadRequestResponse struct {
	Status int    `json:"status" pact:"example=400"`
	Detail string `json:"detail" pact:"example=Enter a percentage between 0 and 100 for lay cases"`
}

type editPaPercentageBadRequestResponse struct {
	Status int    `json:"status" pact:"example=400"`
	Detail string `json:"detail" pact:"example=Enter a percentage between 0 and 100 for PA cases"`
}

type editProPercentageBadRequestResponse struct {
	Status int    `json:"status" pact:"example=400"`
	Detail string `json:"detail" pact:"example=Enter a percentage between 0 and 100 for Pro cases"`
}

const UserExists = "User exists"
const UrlRoute = "/api/v1/random-review-settings"
const BypassMembrane = "OPG-Bypass-Membrane"

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
		name          string
		layPercentage string
		paPercentage  string
		proPercentage string
		reviewCycle   string
		setup         func()
		cookies       []*http.Cookie
		expectedError error
	}{
		{
			name:          "Lay percentage - validation errors",
			layPercentage: "200",
			paPercentage:  "10",
			proPercentage: "18",
			reviewCycle:   "3",
			setup: func() {
				pact.
					AddInteraction().
					Given(UserExists).
					UponReceiving("A request to edit the lay percentage errors on validation").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String(UrlRoute),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN": dsl.String("abcde"),
							"Cookie":       dsl.String("XSRF-TOKEN=abcde; Other=other"),
							BypassMembrane: dsl.String("1"),
						},
						Body: map[string]interface{}{
							"layPercentage": "200",
							"paPercentage":  "10",
							"proPercentage": "18",
							"reviewCycle":   "3",
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusBadRequest,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/problem+json")},
						Body:    dsl.Match(editLayPercentageBadRequestResponse{}),
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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String(UrlRoute),
						Headers: dsl.MapMatcher{
							"Content-type": dsl.String("application/json"),
							"X-XSRF-TOKEN": dsl.String("abcde"),
							"Cookie":       dsl.String("XSRF-TOKEN=abcde; Other=other"),
							BypassMembrane: dsl.String("1"),
						},
						Body: map[string]interface{}{
							"layPercentage": "20",
							"paPercentage":  "10",
							"proPercentage": "18",
							"reviewCycle":   "3",
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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String(UrlRoute),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN": dsl.String("abcde"),
							"Cookie":       dsl.String("XSRF-TOKEN=abcde; Other=other"),
							BypassMembrane: dsl.String("1"),
						},
						Body: map[string]interface{}{
							"layPercentage": "20",
							"paPercentage":  "1000",
							"proPercentage": "18",
							"reviewCycle":   "3",
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusBadRequest,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/problem+json")},
						Body:    dsl.Match(editPaPercentageBadRequestResponse{}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String(UrlRoute),
						Headers: dsl.MapMatcher{
							"Content-type": dsl.String("application/json"),
							"X-XSRF-TOKEN": dsl.String("abcde"),
							"Cookie":       dsl.String("XSRF-TOKEN=abcde; Other=other"),
							BypassMembrane: dsl.String("1"),
						},
						Body: map[string]interface{}{
							"layPercentage": "20",
							"paPercentage":  "50",
							"proPercentage": "18",
							"reviewCycle":   "3",
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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String(UrlRoute),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN": dsl.String("abcde"),
							"Cookie":       dsl.String("XSRF-TOKEN=abcde; Other=other"),
							BypassMembrane: dsl.String("1"),
						},
						Body: map[string]interface{}{
							"layPercentage": "20",
							"paPercentage":  "50",
							"proPercentage": "2000",
							"reviewCycle":   "3",
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusBadRequest,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/problem+json")},
						Body:    dsl.Match(editProPercentageBadRequestResponse{}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String(UrlRoute),
						Headers: dsl.MapMatcher{
							"Content-type": dsl.String("application/json"),
							"X-XSRF-TOKEN": dsl.String("abcde"),
							"Cookie":       dsl.String("XSRF-TOKEN=abcde; Other=other"),
							BypassMembrane: dsl.String("1"),
						},
						Body: map[string]interface{}{
							"layPercentage": "20",
							"paPercentage":  "50",
							"proPercentage": "18",
							"reviewCycle":   "3",
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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String(UrlRoute),
						Headers: dsl.MapMatcher{
							"Content-type": dsl.String("application/json"),
							"X-XSRF-TOKEN": dsl.String("abcde"),
							"Cookie":       dsl.String("XSRF-TOKEN=abcde; Other=other"),
							BypassMembrane: dsl.String("1"),
						},
						Body: map[string]interface{}{
							"layPercentage": "20",
							"paPercentage":  "50",
							"proPercentage": "18",
							"reviewCycle":   "5",
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))
				data := EditRandomReview{tc.layPercentage, tc.paPercentage, tc.proPercentage, tc.reviewCycle}

				err := client.EditRandomReviewSettings(getContext(tc.cookies), data)

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

	err := client.EditRandomReviewSettings(getContext(nil), EditRandomReview{"3", "10", "18", "20"})
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + UrlRoute,
		Method: http.MethodPost,
	}, err)
}
