package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockEditLayReviewCycleClient struct {
	count             int
	saveCount         int
	lastCtx           sirius.Context
    lastRequest       string
	err               error
	data              sirius.RandomReviews
	lastArguments struct {
        LayPercentage   string
        ReviewCycle     string
	}
}

func (m *mockEditLayReviewCycleClient) RandomReviews(ctx sirius.Context) (sirius.RandomReviews, error) {
	m.count += 1
	m.lastCtx = ctx
	m.lastRequest = "RandomReviews"

	return m.data,  m.err
}

func (m *mockEditLayReviewCycleClient) EditLayPercentageReviewCycle(ctx sirius.Context, layPercentage string, reviewCycle string) (error) {
    m.saveCount += 1
	m.lastCtx = ctx
	m.lastRequest = "EditLayPercentageReviewCycle"
    m.lastArguments.LayPercentage = layPercentage
    m.lastArguments.ReviewCycle = reviewCycle

	return m.err
}

func (m *mockEditLayReviewCycleClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-random-review-settings": sirius.PermissionGroup{Permissions: []string{"post"}}}
}

func TestGetLayReviewCycle(t *testing.T) {
	assert := assert.New(t)

    data := sirius.RandomReviews{
     	LayPercentage: 10,
     	ReviewCycle: 1,
    }

	client := &mockEditLayReviewCycleClient{data: data}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := editLayReviewCycle(client, template)
	err := handler(client.requiredPermissions(), w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

    assert.Equal(1, template.count)
    assert.Equal("page", template.lastName)
    assert.Equal(editLayReviewCycleVars{
        Path:        "/path",
        ReviewCycle: "1",
    }, template.lastVars)
}

func TestPostLayReviewCycle(t *testing.T) {
	assert := assert.New(t)

	data := sirius.RandomReviews{
		LayPercentage: 10,
		ReviewCycle: 1,
	}

	client := &mockEditLayReviewCycleClient{data: data}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("layReviewCycle=1"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := editLayReviewCycle(client, template)

	err := handler(client.requiredPermissions(), w, r)
	assert.Equal(Redirect("/random-reviews"), err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(0, template.count)
	assert.Equal(1, client.data.ReviewCycle)
	assert.Equal(10, client.data.LayPercentage)
}

func TestPostLayReviewCycleValidationError(t *testing.T) {

	assert := assert.New(t)

	data := sirius.RandomReviews{
		LayPercentage: 10,
		ReviewCycle: 1,
	}

	errors := sirius.ValidationErrors{
		"x": {
			"y": "z",
		},
	}

	client := &mockEditLayReviewCycleClient{data: data}
	client.err = sirius.ValidationError{
		Errors: errors,
	}

	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("layReviewCycle=test"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := editLayReviewCycle(client, template)

	err := handler(client.requiredPermissions(), w, r)
	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusBadRequest, resp.StatusCode)
	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editLayReviewCycleVars{
	    Path: "/path",
		ReviewCycle: "test",
		Errors: errors,
	}, template.lastVars)

}
