package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockViewTeamClient struct {
	count         int
	lastCookies   []*http.Cookie
	err           error
	data          sirius.Team
	lastRequestID int
}

func (m *mockViewTeamClient) Team(ctx context.Context, cookies []*http.Cookie, id int) (sirius.Team, error) {
	m.count += 1
	m.lastCookies = cookies
	m.lastRequestID = id

	return m.data, m.err
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
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/16", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := viewTeam(client, template, "http://sirius")(w, r)
	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(r.Cookies(), client.lastCookies)

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(viewTeamVars{
		Path:      "/teams/16",
		SiriusURL: "http://sirius",
		Team:      data,
	}, template.lastVars)
}

func TestViewTeamNotFound(t *testing.T) {
	assert := assert.New(t)

	client := &mockViewTeamClient{
		err: StatusError(http.StatusNotFound),
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/25", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := viewTeam(client, template, "http://sirius")(w, r)

	assert.Equal(StatusError(http.StatusNotFound), err)
}

func TestPostViewTeam(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)

	err := viewTeam(nil, nil, "http://sirius")(w, r)
	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
}
