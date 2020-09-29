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

type mockEditMyDetailsClient struct {
	count         int
	saveCount     int
	lastCookies   []*http.Cookie
	lastRequest   string
	err           error
	errSave       error
	data          sirius.MyDetails
	lastArguments struct {
		ID          int
		PhoneNumber string
	}

	permissionsCount int
	hasPermission    bool
	permissionsErr   error
	lastPermission   struct {
		Group  string
		Method string
	}
}

func (m *mockEditMyDetailsClient) MyDetails(ctx context.Context, cookies []*http.Cookie) (sirius.MyDetails, error) {
	m.count += 1
	m.lastCookies = cookies
	m.lastRequest = "MyDetails"

	return m.data, m.err
}

func (m *mockEditMyDetailsClient) EditMyDetails(ctx context.Context, cookies []*http.Cookie, id int, phoneNumber string) error {
	m.saveCount += 1
	m.lastCookies = cookies
	m.lastRequest = "EditMyDetails"
	m.lastArguments.ID = id
	m.lastArguments.PhoneNumber = phoneNumber

	return m.errSave
}

func (m *mockEditMyDetailsClient) HasPermission(ctx context.Context, cookies []*http.Cookie, group string, method string) (bool, error) {
	m.permissionsCount += 1
	m.lastCookies = cookies
	m.lastPermission.Group = group
	m.lastPermission.Method = method

	return m.hasPermission, m.permissionsErr
}

func TestGetEditMyDetails(t *testing.T) {
	assert := assert.New(t)

	data := sirius.MyDetails{
		ID:          123,
		Firstname:   "John",
		Surname:     "Doe",
		Email:       "john@doe.com",
		PhoneNumber: "123",
		Roles:       []string{"A", "COP User", "B"},
		Teams: []sirius.MyDetailsTeam{
			{DisplayName: "A Team"},
		},
	}
	client := &mockEditMyDetailsClient{data: data, hasPermission: true}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	handler := editMyDetails(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(r.Cookies(), client.lastCookies)

	assert.Equal("user", client.lastPermission.Group)
	assert.Equal("patch", client.lastPermission.Method)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editMyDetailsVars{
		Path:        "/path",
		SiriusURL:   "http://sirius",
		PhoneNumber: "123",
	}, template.lastVars)
}

func TestGetEditMyDetailsUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditMyDetailsClient{hasPermission: true, err: sirius.ErrUnauthorized}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := editMyDetails(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Equal(sirius.ErrUnauthorized, err)

	assert.Equal(0, template.count)
}

func TestGetEditMyDetailsNotPermitted(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditMyDetailsClient{hasPermission: false}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := editMyDetails(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Equal(StatusError(http.StatusForbidden), err)

	assert.Equal(0, template.count)
}

func TestGetEditMyDetailsSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditMyDetailsClient{hasPermission: true, err: errors.New("err")}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := editMyDetails(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Equal("err", err.Error())

	assert.Equal(0, template.count)
}

func TestPostEditMyDetails(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditMyDetailsClient{
		hasPermission: true,
		data: sirius.MyDetails{
			ID: 31,
		},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("phonenumber=0189202"))
	r.Header.Add("Content-type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	handler := editMyDetails(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Equal(RedirectError("/my-details"), err)

	assert.Equal(1, client.count)
	assert.Equal(1, client.saveCount)

	assert.Equal(r.Cookies(), client.lastCookies)
	assert.Equal("EditMyDetails", client.lastRequest)
	assert.Equal(31, client.lastArguments.ID)
	assert.Equal("0189202", client.lastArguments.PhoneNumber)

	assert.Equal(0, template.count)
}

func TestPostEditMyDetailsUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditMyDetailsClient{hasPermission: true, errSave: sirius.ErrUnauthorized}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("phonenumber=0189202"))

	handler := editMyDetails(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Equal(sirius.ErrUnauthorized, err)

	assert.Equal(1, client.count)
	assert.Equal(1, client.saveCount)

	assert.Equal(0, template.count)
}

func TestPostEditMyDetailsSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditMyDetailsClient{hasPermission: true, errSave: errors.New("err")}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("phonenumber=0189202"))
	r.Header.Add("Content-type", "application/x-www-form-urlencoded")

	handler := editMyDetails(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Equal("err", err.Error())

	assert.Equal(1, client.count)
	assert.Equal(1, client.saveCount)

	assert.Equal(0, template.count)
}

func TestPostEditMyDetailsInvalidRequest(t *testing.T) {
	assert := assert.New(t)

	validationError := &sirius.ValidationError{
		Errors: sirius.ValidationErrors{
			"phoneNumber": {
				"invalidNumber": "Phone number is not in valid format",
			},
		},
	}

	client := &mockEditMyDetailsClient{hasPermission: true, errSave: validationError}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("phonenumber=invalid+phone+number"))
	r.Header.Add("Content-type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	handler := editMyDetails(nil, client, template, "http://sirius")
	err := handler(w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusBadRequest, resp.StatusCode)

	assert.Equal(1, client.count)
	assert.Equal(1, client.saveCount)

	assert.Equal(r.Cookies(), client.lastCookies)
	assert.Equal("EditMyDetails", client.lastRequest)

	assert.Equal(1, template.count)
	assert.Equal(editMyDetailsVars{
		Path:        "/path",
		SiriusURL:   "http://sirius",
		PhoneNumber: "invalid phone number",
		Errors: map[string]map[string]string{
			"phoneNumber": {
				"invalidNumber": "Phone number is not in valid format",
			},
		},
	}, template.lastVars)
}
