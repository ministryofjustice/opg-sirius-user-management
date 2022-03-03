package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/stretchr/testify/assert"
)

type mockViewTeamClient struct {
	count         int
	lastCtx       sirius.Context
	err           error
	data          sirius.Team
	lastRequestID int
}

func (m *mockViewTeamClient) Team(ctx sirius.Context, id int) (sirius.Team, error) {
	m.count += 1
	m.lastCtx = ctx
	m.lastRequestID = id

	return m.data, m.err
}

func (m *mockViewTeamClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-teams": sirius.PermissionGroup{Permissions: []string{"put"}}}
}

func TestViewTeam(t *testing.T) {
	assert := assert.New(t)

	data := sirius.Team{
		ID:          16,
		DisplayName: "Lay allocations",
		Type:        "Allocations",
		Members: []sirius.TeamMember{
			{
				DisplayName: "Stephani Bennard",
				Email:       "s.bennard@opgtest.com",
			},
		},
	}
	client := &mockViewTeamClient{
		data: data,
	}
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/16", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := viewTeam(client, template.Func())(w, r)
	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal(viewTeamVars{
		Path: "/teams/16",
		Team: data,
	}, template.lastVars)
}

func TestViewTeamNoPermission(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, sirius.PermissionSet{}))

	err := viewTeam(nil, nil)(w, r)
	assert.Equal(handler.Status(http.StatusForbidden), err)
}

func TestViewTeamNotFound(t *testing.T) {
	assert := assert.New(t)

	client := &mockViewTeamClient{
		err: handler.Status(http.StatusNotFound),
	}
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/25", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := viewTeam(client, template.Func())(w, r)

	assert.Equal(handler.Status(http.StatusNotFound), err)
}

func TestViewTeamBadPath(t *testing.T) {
	assert := assert.New(t)

	client := &mockViewTeamClient{}
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/jeoi", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := viewTeam(client, template.Func())(w, r)

	assert.Equal(handler.Status(http.StatusNotFound), err)
}

func TestPostViewTeam(t *testing.T) {
	assert := assert.New(t)

	client := &mockViewTeamClient{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := viewTeam(nil, nil)(w, r)
	assert.Equal(handler.Status(http.StatusMethodNotAllowed), err)
}
