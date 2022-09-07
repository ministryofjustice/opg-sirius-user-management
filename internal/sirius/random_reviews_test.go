package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

const UrlRandomReview = "/api/v1/random-review-settings"

func TestRandomReviews(t *testing.T) {
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
		name                  string
		setup                 func()
		cookies               []*http.Cookie
		expectedRandomReviews RandomReviews
		expectedError         error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get random reviews").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String(UrlRandomReview),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"layPercentage": dsl.Like(20),
							"paPercentage": dsl.Like(30),
							"reviewCycle":   dsl.Like(3),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedRandomReviews: RandomReviews{
				LayPercentage: 20,
				PaPercentage: 30,
				ReviewCycle:   3,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				randomReviews, err := client.RandomReviews(getContext(tc.cookies))
				assert.Equal(t, tc.expectedRandomReviews, randomReviews)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestRandomReviewsStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.RandomReviews(getContext(nil))
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + UrlRandomReview,
		Method: http.MethodGet,
	}, err)
}
