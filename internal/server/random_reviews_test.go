package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockRandomReviewsClient struct {
	count   int
	lastCtx sirius.Context
	err     error
	data    sirius.RandomReviews
}

func (m *mockRandomReviewsClient) RandomReviews(ctx sirius.Context) (sirius.RandomReviews, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.data, m.err
}

func (m *mockRandomReviewsClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-random-review-settings": sirius.PermissionGroup{Permissions: []string{"get"}}}
}

func TestGetRandomReviews(t *testing.T) {
	assert := assert.New(t)

	data := sirius.RandomReviews{
		LayPercentage: 20,
		ReviewCycle:   3,
	}
	client := &mockRandomReviewsClient{data: data}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := randomReviews(client, template)
	err := handler(client.requiredPermissions(), w, r)
	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)
	assert.Equal(1, client.count)
	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(randomReviewsVars{
		Path:          "/path",
		LayPercentage: 20,
		ReviewCycle:   3,
	}, template.lastVars)
}

func TestGetRandomReviewsUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockRandomReviewsClient{err: StatusError(http.StatusForbidden)}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := randomReviews(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal(StatusError(http.StatusForbidden), err)

	assert.Equal(0, template.count)
}

func TestGetRandomReviewsMethodNotAllowed(t *testing.T) {
	assert := assert.New(t)

	client := &mockRandomReviewsClient{err: StatusError(http.StatusMethodNotAllowed)}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)

	handler := randomReviews(client, template)
	err := handler(client.requiredPermissions(), w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, template.count)
}

func TestGetRandomReviewsSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockRandomReviewsClient{err: errors.New("err")}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := randomReviews(client, template)
	err := handler(client.requiredPermissions(), w, r)

	assert.Equal("err", err.Error())

	assert.Equal(0, template.count)
}
