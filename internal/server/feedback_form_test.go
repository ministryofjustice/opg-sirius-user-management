package server

import (
	"errors"
	"github.com/ministryofjustice/opg-sirius-user-management/internal/model"
	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockFeedbackFormClient struct {
	count       int
	lastCtx     sirius.Context
	form        model.FeedbackForm
	addFeedback struct {
		err error
	}
}

func (m *mockFeedbackFormClient) AddFeedback(ctx sirius.Context, form model.FeedbackForm) error {
	m.count += 1
	m.lastCtx = ctx

	return m.addFeedback.err
}

func TestGetFeedbackForm(t *testing.T) {
	assert := assert.New(t)

	client := &mockFeedbackFormClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := feedbackForm(client, template)
	err := handler(sirius.PermissionSet{}, w, r)
	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(1, template.count)
	assert.Equal(feedbackFormVars{
		Path:    "/feedback",
		Success: false,
		Form:    model.FeedbackForm{},
	}, template.lastVars)
}

func TestPostFeedbackForm(t *testing.T) {
	assert := assert.New(t)

	client := &mockFeedbackFormClient{
		form: model.FeedbackForm{
			IsSupervisionFeedback: true,
			Name:                  "",
			Email:                 "",
			CaseNumber:            "",
			Message:               "Im not happy with this service",
		},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/feedback-form", strings.NewReader("more-detail=Im not happy with this service"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := feedbackForm(client, template)(sirius.PermissionSet{}, w, r)
	assert.Nil(err)
	assert.Equal(1, template.count)
	assert.Equal(feedbackFormVars{
		Path:    "/feedback",
		Success: true,
		Form:    model.FeedbackForm{},
	}, template.lastVars)
}

func TestDeputies_MethodNotAllowed(t *testing.T) {
	methods := []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPut,
		http.MethodTrace,
	}
	for _, method := range methods {
		t.Run("Test "+method, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockFeedbackFormClient{
				form: model.FeedbackForm{
					IsSupervisionFeedback: true,
					Name:                  "",
					Email:                 "",
					CaseNumber:            "",
					Message:               "Im not happy with this service",
				},
			}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(method, "/feedback-form", strings.NewReader("more-detail=Im not happy with this service"))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			err := feedbackForm(client, template)(sirius.PermissionSet{}, w, r)
			assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
			assert.Equal(0, template.count)
		})
	}
}

func TestHandlesValidationErrorIfReturnedByAddFeedback(t *testing.T) {
	assert := assert.New(t)

	validationErrors := sirius.ValidationErrors{
		"feedback": {
			"stringLengthTooLong": "Message is too long",
		},
	}

	client := &mockFeedbackFormClient{
		form: model.FeedbackForm{
			Message: "test",
		},
	}
	client.addFeedback.err = sirius.ValidationError{
		Errors: validationErrors,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("more-detail=test"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := feedbackForm(client, template)(sirius.PermissionSet{}, w, r)
	assert.Nil(err)

	assert.Equal(1, template.count)
	assert.Equal(feedbackFormVars{
		Path:    "/feedback",
		Success: false,
		Errors:  validationErrors,
		Form:    model.FeedbackForm{Message: "test", IsSupervisionFeedback: true},
	}, template.lastVars)
	assert.Equal(1, client.count)
}

func TestAddFeedbackFormError(t *testing.T) {
	assert := assert.New(t)
	expectedError := sirius.ClientError("problem")

	client := &mockFeedbackFormClient{
		form: model.FeedbackForm{
			IsSupervisionFeedback: true,
			Name:                  "",
			Email:                 "",
			CaseNumber:            "",
			Message:               "Im not happy with this service",
		},
	}
	client.addFeedback.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/feedback-form", strings.NewReader("more-detail=Im not happy with this service"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := feedbackForm(client, template)(sirius.PermissionSet{}, w, r)
	assert.Equal(expectedError, err)
	assert.Equal(1, client.count)
	assert.Equal(0, template.count)
}

func TestHandlesErrorIfReturned(t *testing.T) {
	assert := assert.New(t)
	expectedError := errors.New("oops")

	client := &mockFeedbackFormClient{
		form: model.FeedbackForm{
			Message: "test",
		},
	}
	client.addFeedback.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/feedback-form", strings.NewReader("more-detail=test"))

	handler := feedbackForm(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal(expectedError, err)
	assert.Equal(1, client.count)
	assert.Equal(0, template.count)
}
