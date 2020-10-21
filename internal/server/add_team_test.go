package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockAddTeamClient struct {
	addTeam struct {
		count        int
		lastCookies  []*http.Cookie
		lastName     string
		lastTeamType string
		lastPhone    string
		lastEmail    string
		data         int
		err          error
	}
	teamTypes struct {
		count       int
		lastCookies []*http.Cookie
		data        []sirius.RefDataTeamType
		err         error
	}
}

func (m *mockAddTeamClient) AddTeam(ctx context.Context, cookies []*http.Cookie, name, teamType, phone, email string) (int, error) {
	m.addTeam.count += 1
	m.addTeam.lastCookies = cookies
	m.addTeam.lastName = name
	m.addTeam.lastTeamType = teamType
	m.addTeam.lastPhone = phone
	m.addTeam.lastEmail = email

	return m.addTeam.data, m.addTeam.err
}

func (m *mockAddTeamClient) TeamTypes(ctx context.Context, cookies []*http.Cookie) ([]sirius.RefDataTeamType, error) {
	m.teamTypes.count += 1
	m.teamTypes.lastCookies = cookies

	return m.teamTypes.data, m.teamTypes.err
}

func TestGetAddTeam(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddTeamClient{}
	client.teamTypes.data = []sirius.RefDataTeamType{
		{Handle: "a"},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := addTeam(client, template, "http://sirius")(w, r)
	assert.Nil(err)

	assert.Equal(0, client.addTeam.count)

	assert.Equal(1, client.teamTypes.count)
	assert.Equal(r.Cookies(), client.teamTypes.lastCookies)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addTeamVars{
		Path:      "/path",
		SiriusURL: "http://sirius",
		TeamTypes: client.teamTypes.data,
	}, template.lastVars)
}

func TestGetAddTeamTeamTypesError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockAddTeamClient{}
	client.teamTypes.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := addTeam(client, template, "http://sirius")(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.teamTypes.count)
	assert.Equal(0, client.addTeam.count)
	assert.Equal(0, template.count)
}

func TestPostAddTeam(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddTeamClient{}
	client.addTeam.data = 123
	client.teamTypes.data = []sirius.RefDataTeamType{
		{Handle: "a"},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("name=a&service=b&supervision-type=c&phone=d&email=e"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := addTeam(client, template, "http://sirius")(w, r)
	assert.Equal(RedirectError("/teams/123"), err)

	assert.Equal(1, client.addTeam.count)
	assert.Equal(r.Cookies(), client.addTeam.lastCookies)
	assert.Equal("a", client.addTeam.lastName)
	assert.Equal("c", client.addTeam.lastTeamType)
	assert.Equal("d", client.addTeam.lastPhone)
	assert.Equal("e", client.addTeam.lastEmail)

	assert.Equal(0, client.teamTypes.count)
	assert.Equal(0, template.count)
}

func TestPostAddTeamLpa(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddTeamClient{}
	client.addTeam.data = 123
	client.teamTypes.data = []sirius.RefDataTeamType{
		{Handle: "a"},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("name=a&service=lpa&supervision-type=c&phone=d&email=e"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := addTeam(client, template, "http://sirius")(w, r)
	assert.Equal(RedirectError("/teams/123"), err)

	assert.Equal(1, client.addTeam.count)
	assert.Equal(r.Cookies(), client.addTeam.lastCookies)
	assert.Equal("a", client.addTeam.lastName)
	assert.Equal("", client.addTeam.lastTeamType)
	assert.Equal("d", client.addTeam.lastPhone)
	assert.Equal("e", client.addTeam.lastEmail)

	assert.Equal(0, client.teamTypes.count)
	assert.Equal(0, template.count)
}

func TestPostAddTeamValidationError(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddTeamClient{}
	client.teamTypes.data = []sirius.RefDataTeamType{
		{Handle: "a"},
	}
	client.addTeam.err = sirius.ValidationError{
		Errors: sirius.ValidationErrors{
			"something": {"": "something"},
		},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("name=a&service=b&supervision-type=c&phone=d&email=e"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := addTeam(client, template, "http://sirius")(w, r)
	assert.Nil(err)

	assert.Equal(1, client.addTeam.count)
	assert.Equal(1, client.teamTypes.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addTeamVars{
		Path:      "/path",
		SiriusURL: "http://sirius",
		Name:      "a",
		Service:   "b",
		TeamType:  "c",
		Phone:     "d",
		Email:     "e",
		TeamTypes: client.teamTypes.data,
		Errors: sirius.ValidationErrors{
			"something": {"": "something"},
		},
	}, template.lastVars)
}

func TestPostAddTeamError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockAddTeamClient{}
	client.teamTypes.data = []sirius.RefDataTeamType{
		{Handle: "a"},
	}
	client.addTeam.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("name=a&service=b&supervision-type=c&phone=d&email=e"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	err := addTeam(client, template, "http://sirius")(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.addTeam.count)
	assert.Equal(0, client.teamTypes.count)
	assert.Equal(0, template.count)
}
func TestPutAddTeam(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddTeamClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/path", nil)

	err := addTeam(client, nil, "http://sirius")(w, r)
	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
}
