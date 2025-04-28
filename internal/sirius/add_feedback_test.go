package sirius

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-user-management/internal/model"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestAddFeedback(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		form          model.FeedbackForm
		expectedError func(int) error
	}{
		{
			name: "OK",
			form: model.FeedbackForm{
				IsSupervisionFeedback: true,
				Message:               "some feedback",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("Supervision team with members exists").
					UponReceiving("A request to add feedback").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String(SupervisionAPIPath + "/supervision-feedback"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"isSupervisionFeedback": true,
							"name":                  "",
							"email":                 "",
							"caseNumber":            "",
							"message":               "some feedback",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusForbidden,
					})
			},
			expectedError: func(port int) error {
				return StatusError{Code: 403, URL: fmt.Sprintf("http://127.0.0.1:%d/supervision-api/supervision-feedback", port), Method: http.MethodPost}
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.AddFeedback(Context{Context: context.Background()}, tc.form)

				assert.Equal(t, tc.expectedError(config.Port), err)
				return nil
			}))
		})
	}
}

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

func TestAddFeedbackCanHandleBadRequest(t *testing.T) {
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
		URL:    svr.URL + SupervisionAPIPath + "/supervision-feedback",
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
		URL:    svr.URL + SupervisionAPIPath + "/supervision-feedback",
		Method: http.MethodPost,
	}, err)
}

func TestCanThrowReqErrors(t *testing.T) {
	jsonResponse := `{"detail": "Could not post to Slack"}`

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	r := io.NopCloser(bytes.NewReader([]byte(jsonResponse)))

	mocks.GetDoFunc = func(rq *http.Request) (*http.Response, error) {

		err := json.NewDecoder(rq.Body)
		assert.Nil(t, err)

		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	err := client.AddFeedback(getContextForMock(nil), model.FeedbackForm{
		IsSupervisionFeedback: true,
		Name:                  "",
		Email:                 "",
		CaseNumber:            "",
		Message:               "feedback message",
	})
	assert.Equal(t, nil, err)
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
	}, err)
}

func TestValidationErrorUnmarshalled(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	jsonResponse := `{
    "validation_errors": {
        "feedback": {
            "isEmpty": "Message is required and can not be empty"
        }
    },
    "type": "http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html",
    "title": "Bad Request",
    "status": 400,
    "detail": "Payload failed validation"
	}`
	r := io.NopCloser(bytes.NewReader([]byte(jsonResponse)))

	mocks.GetDoFunc = func(rq *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Body:       r,
		}, nil
	}

	expectedResponse := ValidationError{Message: "", Errors: ValidationErrors{"feedback": map[string]string{"isEmpty": "Message is required and can not be empty"}}}

	err := client.AddFeedback(getContextForMock(nil), model.FeedbackForm{
		IsSupervisionFeedback: true,
		Name:                  "",
		Email:                 "",
		CaseNumber:            "",
		Message:               "feedback message",
	})
	assert.Equal(t, expectedResponse, err)
}

func getContextForMock(cookies []*http.Cookie) Context {
	return Context{
		Context:   context.Background(),
		Cookies:   cookies,
		XSRFToken: "abcde",
	}
}
