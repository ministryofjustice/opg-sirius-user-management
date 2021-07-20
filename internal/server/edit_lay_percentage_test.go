package server

import (
	"net/http"
	"net/http/httptest"
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
	errSave           error
	data              sirius.RandomReviews
	lastArguments struct {
        LayPercentage   string
        ReviewCycle     int
	}
}

func (m *mockEditLayPercentageClient) RandomReviews(ctx sirius.Context) (sirius.RandomReviews, error) {
	m.count += 1
	m.lastCtx = ctx
	m.lastRequest = "RandomReviews"

	return m.data,  m.err
}

func (m *mockEditLayPercentageClient) EditLayPercentage(ctx sirius.Context, layPercentage string, reviewCycle int) (error) {
    m.saveCount += 1
	m.lastCtx = ctx
	m.lastRequest = "EditLayPercentage"
    m.lastArguments.LayPercentage = layPercentage
    m.lastArguments.ReviewCycle = reviewCycle

	return m.errSave
}

func (m *mockEditLayPercentageClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-random-review-settings": sirius.PermissionGroup{Permissions: []string{"post"}}}
}

func TestLayPercentageClient(t *testing.T) {
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
	err := handler(sirius.PermissionSet{}, w, r)

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
