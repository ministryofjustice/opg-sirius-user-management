package sirius

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type ValidationErrors map[string]map[string]string

type MyDetails struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	PhoneNumber string          `json:"phoneNumber"`
	Teams       []MyDetailsTeam `json:"teams"`
	DisplayName string          `json:"displayName"`
	Deleted     bool            `json:"deleted"`
	Email       string          `json:"email"`
	Firstname   string          `json:"firstname"`
	Surname     string          `json:"surname"`
	Roles       []string        `json:"roles"`
	Locked      bool            `json:"locked"`
	Suspended   bool            `json:"suspended"`
}

type MyDetailsTeam struct {
	DisplayName string `json:"displayName"`
}

func (c *Client) MyDetails(ctx context.Context, cookies []*http.Cookie) (MyDetails, error) {
	var v MyDetails

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/users/current", nil, cookies)
	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, errors.New("returned non-2XX response: " + strconv.Itoa(resp.StatusCode))
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}

func (c *Client) EditMyDetails(ctx context.Context, cookies []*http.Cookie, id int, phoneNumber string) (ValidationErrors, error) {
	var v struct {
		Status           int              `json:"status"`
		Detail           string           `json:"detail"`
		ValidationErrors ValidationErrors `json:"validation_errors"`
	}

	var body = strings.NewReader("{\"phoneNumber\":\"" + phoneNumber + "\"}")

	req, err := c.newRequest(
		ctx,
		http.MethodPut,
		"/api/v1/users/"+strconv.Itoa(id)+"/updateTelephoneNumber",
		body,
		cookies,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		json.NewDecoder(resp.Body).Decode(&v)
		return v.ValidationErrors, nil
	}

	return nil, nil
}
