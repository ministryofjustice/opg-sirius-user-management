package sirius

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type UserStatus string

func (us UserStatus) TagColour() string {
	if us == "Suspended" {
		return "govuk-tag--grey"
	} else if us == "Locked" {
		return "govuk-tag--orange"
	} else {
		return ""
	}
}

type User struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Locked      bool   `json:"locked"`
	Suspended   bool   `json:"suspended"`
	Status      UserStatus
}

func (c *Client) ListUsers(ctx context.Context, cookies []*http.Cookie) ([]User, error) {
	var v []User

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

	for key, user := range v {
		if user.Suspended {
			v[key].Status = "Suspended"
		} else if user.Locked {
			v[key].Status = "Locked"
		} else {
			v[key].Status = "Active"
		}
	}

	return v, nil
}
