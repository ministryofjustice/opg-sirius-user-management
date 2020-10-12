package sirius

import (
	"context"
	"encoding/json"
	"net/http"
)

type apiTeam struct {
	ID          int        `json:"id"`
	DisplayName string     `json:"displayName"`
	Members     []struct{} `json:"members"`
	TeamType    *struct {
		Label string `json:"label"`
	} `json:"teamType"`
}

type Team struct {
	ID          int
	DisplayName string
	Members     int
	Type        string
}

func (c *Client) Teams(ctx context.Context, cookies []*http.Cookie) ([]Team, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/teams", nil, cookies)
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
		return nil, newStatusError(resp)
	}

	var v []apiTeam
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	teams := make([]Team, len(v))

	for i, t := range v {
		teams[i] = Team{
			ID:          t.ID,
			DisplayName: t.DisplayName,
			Members:     len(t.Members),
			Type:        "LPA",
		}

		if t.TeamType != nil {
			teams[i].Type = "Supervision — " + t.TeamType.Label
		}
	}

	return teams, nil
}
