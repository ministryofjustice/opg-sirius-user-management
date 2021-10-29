package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockEditRandomReviewSettingsClient struct {
	count         int
	saveCount     int
	lastCtx       sirius.Context
	lastRequest   string
	err           error
	data          sirius.RandomReviews
	lastArguments struct {
		layPercentage string
		paPercentage  string
		reviewCycle   string
	}
}

func (m *mockEditRandomReviewSettingsClient) RandomReviews(ctx sirius.Context) (sirius.RandomReviews, error) {
	m.count += 1
	m.lastCtx = ctx
	m.lastRequest = "RandomReviews"

	return m.data, m.err
}

func (m *mockEditRandomReviewSettingsClient) EditRandomReviewSettings(ctx sirius.Context, reviewSettings sirius.EditRandomReview) error {
	m.saveCount += 1
	m.lastCtx = ctx
	m.lastRequest = "EditRandomReviewSettings"
	m.lastArguments.layPercentage = reviewSettings.LayPercentage
	m.lastArguments.paPercentage = reviewSettings.PaPercentage
	m.lastArguments.reviewCycle = reviewSettings.ReviewCycle

	return m.err
}

func (m *mockEditRandomReviewSettingsClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-random-review-settings": sirius.PermissionGroup{Permissions: []string{"post"}}}
}

func TestGetRandomReviewSettings(t *testing.T) {
	assert := assert.New(t)

	data := sirius.RandomReviews{
		LayPercentage: 10,
		PaPercentage: 20,
		ReviewCycle:   1,
	}

	client := &mockEditRandomReviewSettingsClient{data: data}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := editRandomReviewSettings(client, template)
	err := handler(client.requiredPermissions(), w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editRandomReviewSettingsVars{
		Path:        "/path",
		LayPercentage: 10,
		PaPercentage: 20,
		ReviewCycle:   1,
	}, template.lastVars)
}

func TestPostRandomReviewSettings(t *testing.T) {
	assert := assert.New(t)

	data := sirius.RandomReviews{
		LayPercentage: 10,
		PaPercentage: 20,
		ReviewCycle:   1,
	}

	client := &mockEditRandomReviewSettingsClient{data: data}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("layReviewCycle=1"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := editRandomReviewSettings(client, template)

	err := handler(client.requiredPermissions(), w, r)
	assert.Equal(RedirectError("/random-reviews"), err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(0, template.count)
	assert.Equal(1, client.data.ReviewCycle)
	assert.Equal(10, client.data.LayPercentage)
	assert.Equal(20, client.data.PaPercentage)
}

func TestPostRandomReviewSettingsValidationError(t *testing.T) {

	assert := assert.New(t)

	data := sirius.RandomReviews{
		LayPercentage: 10,
		PaPercentage: 20,
		ReviewCycle:   1,
	}

	errors := sirius.ValidationErrors{
		"x": {
			"y": "z",
		},
	}

	client := &mockEditRandomReviewSettingsClient{data: data}
	client.err = sirius.ValidationError{
		Errors: errors,
	}

	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("reviewCycle=-1"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := editRandomReviewSettings(client, template)

	err := handler(client.requiredPermissions(), w, r)
	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusBadRequest, resp.StatusCode)
	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editRandomReviewSettingsVars{
		Path:        "/path",
		LayPercentage: 10,
		PaPercentage: 20,
		ReviewCycle:   1,
		Errors:      errors,
	}, template.lastVars)
}
