package sirius

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type UserStatus string

func (us UserStatus) String() string {
	return string(us)
}

func (us UserStatus) TagColour() string {
	if us == "Suspended" {
		return "govuk-tag--grey"
	} else {
		return ""
	}
}

type apiUser struct {
	ID          int       `json:"id"`
	DisplayName string    `json:"displayName"`
	Surname     string    `json:"surname"`
	Email       string    `json:"email"`
	Suspended   bool      `json:"suspended"`
	Teams       []apiTeam `json:"teams"`
}

type User struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Status      UserStatus
	Team        string `json:"team"`
}

func (c *Client) SearchUsers(ctx Context, search string) ([]User, error) {
	if len(search) < 3 {
		return nil, ClientError("Search term must be at least three characters")
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/search/users?includeSuspended=1&query="+url.QueryEscape(search), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck // no need to check error when closing body

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v []apiUser
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	sort.SliceStable(v, func(i, j int) bool {
		if strings.EqualFold(v[i].Surname, v[j].Surname) {
			return strings.ToLower(v[i].DisplayName) < strings.ToLower(v[j].DisplayName)
		}

		return strings.ToLower(v[i].Surname) < strings.ToLower(v[j].Surname)
	})

	var users []User
	for _, u := range v {
		var teamName string
		if len(u.Teams) > 0 {
			teamName = u.Teams[0].DisplayName
		}

		user := User{
			ID:          u.ID,
			DisplayName: u.DisplayName,
			Email:       u.Email,
			Status:      "Active",
			Team:        teamName,
		}

		if u.Suspended {
			user.Status = "Suspended"
		}

		users = append(users, user)
	}

	return users, nil
}
