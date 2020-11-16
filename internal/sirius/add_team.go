package sirius

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) AddTeam(ctx Context, name, teamType, phone, email string) (int, error) {
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

	req, err := c.newRequest(ctx, http.MethodPost, "/api/team", strings.NewReader(form.Encode()))
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

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return 0, ValidationError{Errors: v.Data.ErrorMessages}
		}

		return 0, newStatusError(resp)
	}

	var v apiTeamResponse
	err = json.NewDecoder(resp.Body).Decode(&v)

	return v.Data.ID, err
}
