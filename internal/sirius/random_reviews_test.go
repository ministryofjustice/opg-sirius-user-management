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

const UrlRandomReview = "/supervision-api/v1/random-review-settings"

func TestRandomReviews(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name                  string
		setup                 func()
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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String(UrlRandomReview),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"layPercentage": matchers.Like(20),
							"paPercentage":  matchers.Like(30),
							"reviewCycle":   matchers.Like(3),
						}),
					})
			},
			expectedRandomReviews: RandomReviews{
				LayPercentage: 20,
				PaPercentage:  30,
				ReviewCycle:   3,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				randomReviews, err := client.RandomReviews(Context{Context: context.Background()})
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

	_, err := client.RandomReviews(Context{Context: context.Background()})
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + UrlRandomReview,
		Method: http.MethodGet,
	}, err)
}
