package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/stretchr/testify/assert"
)

type mockAddTeamClient struct {
	addTeam struct {
		count        int
		lastCtx      sirius.Context
		lastName     string
		lastTeamType string
		lastPhone    string
		lastEmail    string
		data         int
		err          error
	}
	teamTypes struct {
		count   int
		lastCtx sirius.Context
		data    []sirius.RefDataTeamType
		err     error
	}
}

func (m *mockAddTeamClient) AddTeam(ctx sirius.Context, name, teamType, phone, email string) (int, error) {
	m.addTeam.count += 1
	m.addTeam.lastCtx = ctx
	m.addTeam.lastName = name
	m.addTeam.lastTeamType = teamType
	m.addTeam.lastPhone = phone
	m.addTeam.lastEmail = email

	return m.addTeam.data, m.addTeam.err
}

func (m *mockAddTeamClient) TeamTypes(ctx sirius.Context) ([]sirius.RefDataTeamType, error) {
	m.teamTypes.count += 1
	m.teamTypes.lastCtx = ctx

	return m.teamTypes.data, m.teamTypes.err
}

func (m *mockAddTeamClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-teams": sirius.PermissionGroup{Permissions: []string{"post"}}}
}

func TestGetAddTeam(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddTeamClient{}
	client.teamTypes.data = []sirius.RefDataTeamType{
		{Handle: "a"},
	}
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := addTeam(client, template.Func())(w, r)
	assert.Nil(err)

	assert.Equal(0, client.addTeam.count)

	assert.Equal(1, client.teamTypes.count)
	assert.Equal(getContext(r), client.teamTypes.lastCtx)

	assert.Equal(1, template.count)
	assert.Equal(addTeamVars{
		Path:      "/path",
		TeamTypes: client.teamTypes.data,
	}, template.lastVars)
}

func TestGetAddTeamNoPermission(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, sirius.PermissionSet{}))

	err := addTeam(nil, nil)(w, r)
	assert.Equal(handler.Status(http.StatusForbidden), err)
}

func TestGetAddTeamTeamTypesError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockAddTeamClient{}
	client.teamTypes.err = expectedError
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := addTeam(client, template.Func())(w, r)
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
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("name=a&service=b&supervision-type=c&phone=d&email=e"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := addTeam(client, template.Func())(w, r)
	assert.Equal(handler.Redirect("/teams/123"), err)

	assert.Equal(1, client.addTeam.count)
	assert.Equal(getContext(r), client.addTeam.lastCtx)
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
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("name=a&service=lpa&supervision-type=c&phone=d&email=e"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := addTeam(client, template.Func())(w, r)
	assert.Equal(handler.Redirect("/teams/123"), err)

	assert.Equal(1, client.addTeam.count)
	assert.Equal(getContext(r), client.addTeam.lastCtx)
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
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("name=a&service=b&supervision-type=c&phone=d&email=e"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := addTeam(client, template.Func())(w, r)
	assert.Nil(err)

	assert.Equal(1, client.addTeam.count)
	assert.Equal(1, client.teamTypes.count)

	assert.Equal(1, template.count)
	assert.Equal(addTeamVars{
		Path:      "/path",
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
	template := &mockTemplateFn{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("name=a&service=b&supervision-type=c&phone=d&email=e"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := addTeam(client, template.Func())(w, r)
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
	r = r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, client.requiredPermissions()))

	err := addTeam(client, nil)(w, r)
	assert.Equal(handler.Status(http.StatusMethodNotAllowed), err)
}
