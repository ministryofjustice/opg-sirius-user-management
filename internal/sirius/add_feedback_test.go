package sirius

import (
	"context"
	"github.com/ministryofjustice/opg-sirius-user-management/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//func TestAddFeedback(t *testing.T) {
//	pact := &dsl.Pact{
//		Consumer:          "sirius-user-management",
//		Provider:          "sirius",
//		Host:              "localhost",
//		PactFileWriteMode: "merge",
//		LogDir:            "../../logs",
//		PactDir:           "../../pacts",
//	}
//	defer pact.Teardown()
//
//	testCases := []struct {
//		scenario      string
//		setup         func()
//		name          string
//		message       string
//		expectedError error
//	}{
//		{
//			scenario: "Created",
//			setup: func() {
//				pact.
//					AddInteraction().
//					Given("An admin user").
//					UponReceiving("A feedback request").
//					WithRequest(dsl.Request{
//						Method: http.MethodPost,
//						Path:   dsl.String("/api/supervision-feedback"),
//						Headers: dsl.MapMatcher{
//							"Content-Type": dsl.String("application/json"),
//						},
//						Body: map[string]interface{}{
//							"name":    "testteam",
//							"message": "feedback message",
//						},
//					}).
//					WillRespondWith(dsl.Response{
//						Status: http.StatusInternalServerError,
//					})
//			},
//			name:    "testteam",
//			message: "feedback message",
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.scenario, func(t *testing.T) {
//			tc.setup()
//
//			assert.Nil(t, pact.Verify(func() error {
//				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))
//
//				err := client.AddFeedback(Context{Context: context.Background()}, model.FeedbackForm{Message: tc.message, Name: tc.name})
//				assert.Equal(t, tc.expectedError, err)
//				return nil
//			}))
//		})
//	}
//}

func TestAddFeedbackCanPost(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.AddFeedback(Context{Context: context.Background()}, model.FeedbackForm{
		IsSupervisionFeedback: true,
		Name:                  "Toad",
		Email:                 "toad@toadhall.com",
		CaseNumber:            "123",
		Message:               "content",
	})
	assert.Nil(t, err)
}

func TestAddFeedbackCanHandleValidationError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.AddFeedback(Context{Context: context.Background()}, model.FeedbackForm{
		IsSupervisionFeedback: true,
		Name:                  "Toad",
		Email:                 "toad@toadhall.com",
		CaseNumber:            "123",
		Message:               "content",
	})
	assert.Equal(t, StatusError{
		Code:   http.StatusBadRequest,
		URL:    svr.URL + "/api/supervision-feedback",
		Method: http.MethodPost,
	}, err)
}

func TestAddFeedbackCanHandleUnauthorizedError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.AddFeedback(Context{Context: context.Background()}, model.FeedbackForm{
		IsSupervisionFeedback: true,
		Name:                  "Toad",
		Email:                 "toad@toadhall.com",
		CaseNumber:            "123",
		Message:               "content",
	})
	assert.Equal(t, ClientError("unauthorized"), err)
}

func TestGetCaseloadListCanThrow500Error(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.AddFeedback(Context{Context: context.Background()}, model.FeedbackForm{
		IsSupervisionFeedback: true,
		Name:                  "Toad",
		Email:                 "toad@toadhall.com",
		CaseNumber:            "123",
		Message:               "feedback",
	})
	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/api/supervision-feedback",
		Method: http.MethodPost,
	}, err)
}

func TestAddFeedbackIsEmptyValidationError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

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
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

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
