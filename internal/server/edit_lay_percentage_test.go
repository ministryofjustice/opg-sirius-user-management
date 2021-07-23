package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockEditLayPercentageClient struct {
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

func (m *mockEditLayPercentageClient) RandomReviews(ctx sirius.Context) (sirius.RandomReviews, error) {
	m.count += 1
	m.lastCtx = ctx
	m.lastRequest = "RandomReviews"

	return m.data,  m.err
}

func (m *mockEditLayPercentageClient) EditLayPercentageReviewCycle(ctx sirius.Context, layPercentage string, reviewCycle string) (error) {
    m.saveCount += 1
	m.lastCtx = ctx
	m.lastRequest = "EditLayPercentageReviewCycle"
    m.lastArguments.LayPercentage = layPercentage
    m.lastArguments.ReviewCycle = reviewCycle

	return m.err
}

func (m *mockEditLayPercentageClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-random-review-settings": sirius.PermissionGroup{Permissions: []string{"post"}}}
}

func TestGetLayPercentage(t *testing.T) {
	assert := assert.New(t)

    data := sirius.RandomReviews{
     	LayPercentage: 10,
     	ReviewCycle: 1,
    }

	client := &mockEditLayPercentageClient{data: data}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := editLayPercentage(client, template)
	err := handler(client.requiredPermissions(), w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

    assert.Equal(1, template.count)
    assert.Equal("page", template.lastName)
    assert.Equal(editLayPercentageVars{
        Path:        "/path",
        LayPercentage: "10",
    }, template.lastVars)
}

func TestPostLayPercentage(t *testing.T) {
	assert := assert.New(t)

    data := sirius.RandomReviews{
     	LayPercentage: 10,
     	ReviewCycle: 1,
    }

	client := &mockEditLayPercentageClient{data: data}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("layPercentage=10"))
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := editLayPercentage(client, template)

	err := handler(client.requiredPermissions(), w, r)
	assert.Equal(Redirect("/random-reviews"), err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

    assert.Equal(0, template.count)
    assert.Equal(10, client.data.LayPercentage)
    assert.Equal(1, client.data.ReviewCycle)
}

func TestPostLayPercentageValidationError(t *testing.T) {

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

	client := &mockEditLayPercentageClient{data: data}
	client.err = sirius.ValidationError{
		Errors: errors,
	}

	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("layPercentage=test"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := editLayPercentage(client, template)

	err := handler(client.requiredPermissions(), w, r)
	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusBadRequest, resp.StatusCode)
	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editLayPercentageVars{
	    Path: "/path",
		LayPercentage: "test",
		Errors: errors,
	}, template.lastVars)

}

func TestPostEditLayPercentageNoPermission(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := editLayPercentage(nil, nil)(sirius.PermissionSet{}, w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)
}

func TestGetRandomReviewError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockEditLayPercentageClient{}
	client.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/random-review-settings", nil)

	err := editLayPercentage(client, template)(client.requiredPermissions(), w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.count)
	assert.Equal(0, template.count)
}

func TestPostEditLayPercentageError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockEditLayPercentageClient{}
	client.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("layPercentage=test"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := editLayPercentage(client, template)(client.requiredPermissions(), w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.count)
	assert.Equal(0, template.count)
}

func TestPutEditLayPercentageError(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditLayPercentageClient{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/random-review-settings", nil)

	err := editLayPercentage(nil, nil)(client.requiredPermissions(), w, r)
	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
}
