package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockListTeamsClient struct {
	count       int
	lastCookies []*http.Cookie
	err         error
	data        []sirius.Team
}

func (m *mockListTeamsClient) Teams(ctx context.Context, cookies []*http.Cookie) ([]sirius.Team, error) {
	m.count += 1
	m.lastCookies = cookies

	return m.data, m.err
}

func TestListTeams(t *testing.T) {
	assert := assert.New(t)

	data := []sirius.Team{
		{
			ID:          29,
			DisplayName: "Milo Nihei",
			Members:     make([]sirius.TeamMember, 10),
			Type:        "Top Notch",
		},
	}
	client := &mockListTeamsClient{
		data: data,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := listTeams(client, template, "http://sirius")(w, r)
	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(r.Cookies(), client.lastCookies)

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(listTeamsVars{
		Path:      "/path",
		SiriusURL: "http://sirius",
		Teams:     data,
	}, template.lastVars)
}

func TestListTeamsSearch(t *testing.T) {
	assert := assert.New(t)

	data := []sirius.Team{
		{
			ID:          29,
			DisplayName: "Milo Nihei",
			Members:     make([]sirius.TeamMember, 10),
			Type:        "Top Notch",
		},
		{
			ID:          3,
			DisplayName: "Who",
			Members:     make([]sirius.TeamMember, 5),
			Type:        "Terrible",
		},
	}
	client := &mockListTeamsClient{
		data: data,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?search=milo", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := listTeams(client, template, "http://sirius")(w, r)
	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(r.Cookies(), client.lastCookies)

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(listTeamsVars{
		Path:      "/path",
		SiriusURL: "http://sirius",
		Search:    "milo",
		Teams: []sirius.Team{
			{
				ID:          29,
				DisplayName: "Milo Nihei",
				Members:     make([]sirius.TeamMember, 10),
				Type:        "Top Notch",
			},
		},
	}, template.lastVars)
}

func TestListTeamsError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("err")
	client := &mockListTeamsClient{err: expectedErr}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/?search=long", nil)

	err := listTeams(client, template, "http://sirius")(w, r)

	assert.Equal(expectedErr, err)
	assert.Equal(0, template.count)
}

func TestPostListTeams(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)

	err := listTeams(nil, nil, "http://sirius")(w, r)
	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
}
