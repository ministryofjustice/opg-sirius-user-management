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

func TestAddFeedbackReturns403AsNoWebHookLocally(t *testing.T) {
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
		setup         func()
		userID        int
		cookies       []*http.Cookie
		expectedError error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("A user").
					UponReceiving("A feedback request").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/api/supervision-feedback"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusForbidden,
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.AddFeedback(Context{Context: context.Background()}, model.FeedbackForm{
					IsSupervisionFeedback: true,
					Name:                  "toad",
					Email:                 "toad@toad.com",
					CaseNumber:            "123",
					Message:               "message here",
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
