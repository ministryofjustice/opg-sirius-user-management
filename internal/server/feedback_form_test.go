package server

import (
	"github.com/ministryofjustice/opg-sirius-user-management/internal/model"
	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockFeedbackFormClient struct {
	count   int
	lastCtx sirius.Context
	err     error
	form    model.FeedbackForm
}

func (m *mockFeedbackFormClient) AddFeedback(ctx sirius.Context, form model.FeedbackForm) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
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
		Path: "/feedback",
	}, template.lastVars)
}

func TestConfirmPostFeedbackForm(t *testing.T) {
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

func TestGetFeedbackFormMethodNotAllowed(t *testing.T) {
	assert := assert.New(t)

	client := &mockFeedbackFormClient{
		err: StatusError(http.StatusMethodNotAllowed),
		form: model.FeedbackForm{
			Message: "test",
		},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/feedback-form", strings.NewReader("more-detail=Im not happy with this service"))

	handler := feedbackForm(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
	assert.Equal(0, template.count)
}
