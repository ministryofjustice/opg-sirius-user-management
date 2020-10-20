package sirius

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) AddTeam(ctx context.Context, cookies []*http.Cookie, name, teamType, phone, email string) (int, error) {
	form := url.Values{
		"name":     {name},
		"phone":    {phone},
		"email":    {email},
		"type":     {""},
		"teamType": {""},
	}

	if teamType != "" {
		form.Add("teamType[handle]", teamType)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/team", strings.NewReader(form.Encode()), cookies)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return 0, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusCreated {
		var v struct {
			Data struct {
				ErrorMessages ValidationErrors `json:"errorMessages"`
			} `json:"data"`
		}

		err = json.NewDecoder(resp.Body).Decode(&v)
		if err == nil {
			return 0, ValidationError{Errors: v.Data.ErrorMessages}
		}

		if err == io.EOF {
			return 0, newStatusError(resp)
		}

		return 0, err
	}

	var v apiTeamResponse
	err = json.NewDecoder(resp.Body).Decode(&v)

	return v.Data.ID, err
}
