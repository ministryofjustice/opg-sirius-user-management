package sirius

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type UserStatus string

func (us UserStatus) String() string {
	return string(us)
}

func (us UserStatus) TagColour() string {
	if us == "Suspended" {
		return "govuk-tag--grey"
	} else if us == "Locked" {
		return "govuk-tag--orange"
	} else {
		return ""
	}
}

type apiUser struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	Locked      bool   `json:"locked"`
	Suspended   bool   `json:"suspended"`
}

type User struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Status      UserStatus
}

func (c *Client) ListUsers(ctx context.Context, cookies []*http.Cookie) ([]User, error) {
	var v []apiUser

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/users", nil, cookies)
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
		return nil, errors.New("returned non-2XX response: " + strconv.Itoa(resp.StatusCode))
	}

	err = json.NewDecoder(resp.Body).Decode(&v)

	if err != nil {
		return nil, err
	}

	var users []User

	for _, u := range v {
		user := User{
			ID:          u.ID,
			DisplayName: u.DisplayName,
			Email:       u.Email,
			Status:      "Active",
		}

		if u.Suspended {
			user.Status = "Suspended"
		} else if u.Locked {
			user.Status = "Locked"
		}

		users = append(users, user)
	}

	sort.SliceStable(users, func(i, j int) bool {
		return strings.ToLower(v[i].Surname) < strings.ToLower(v[j].Surname)
	})

	return users, nil
}
