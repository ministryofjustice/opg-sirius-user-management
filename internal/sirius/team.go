package sirius

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type apiTeamResponse struct {
	Data apiTeam `json:"data"`
}

func (c *Client) Team(ctx context.Context, cookies []*http.Cookie, id int) (Team, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/team/"+strconv.Itoa(id), nil, cookies)
	if err != nil {
		return Team{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return Team{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return Team{}, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return Team{}, newStatusError(resp)
	}

	var v apiTeamResponse
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return Team{}, err
	}

	team := Team{
		ID:          v.Data.ID,
		DisplayName: v.Data.DisplayName,
		Type:        "LPA",
	}

	for _, m := range v.Data.Members {
		team.Members = append(team.Members, TeamMember{
			DisplayName: m.DisplayName,
			Email:       m.Email,
		})
	}

	if v.Data.TeamType != nil {
		team.Type = "Supervision — " + v.Data.TeamType.Label
	}

	return team, nil
}
