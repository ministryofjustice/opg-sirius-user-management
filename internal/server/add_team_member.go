package server

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type AddTeamMemberClient interface {
	Team(context.Context, []*http.Cookie, int) (sirius.Team, error)
	EditTeam(context.Context, []*http.Cookie, sirius.Team) error
	SearchUsers(context.Context, []*http.Cookie, string) ([]sirius.User, error)
}

type addTeamMemberVars struct {
	Path      string
	SiriusURL string
	Search    string
	Team      sirius.Team
	Users     []sirius.User
	Members   map[int]bool
	Success   string
	Errors    sirius.ValidationErrors
}

func addTeamMember(client AddTeamMemberClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/teams/add-member/"))
		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		team, err := client.Team(r.Context(), r.Cookies(), id)
		if err != nil {
			return err
		}

		vars := addTeamMemberVars{
			Path:      r.URL.Path,
			SiriusURL: siriusURL,
			Team:      team,
		}

		if r.Method == http.MethodPost {
			memberID, err := strconv.Atoi(r.PostFormValue("id"))
			if err != nil {
				return err
			}

			team.Members = append(team.Members, sirius.TeamMember{ID: memberID})

			err = client.EditTeam(r.Context(), r.Cookies(), team)

			if _, ok := err.(sirius.ClientError); ok {
				vars.Errors = sirius.ValidationErrors{
					"search": {
						"": err.Error(),
					},
				}
				w.WriteHeader(http.StatusBadRequest)
			} else if err != nil {
				return err
			} else {
				vars.Success = r.PostFormValue("email")
			}
		}

		vars.Search = r.FormValue("search")

		if vars.Search != "" {
			users, err := client.SearchUsers(r.Context(), r.Cookies(), vars.Search)

			if _, ok := err.(sirius.ClientError); ok {
				vars.Errors = sirius.ValidationErrors{
					"search": {
						"": err.Error(),
					},
				}
			} else if err != nil {
				return err
			} else {
				members := map[int]bool{}
				for _, member := range team.Members {
					members[member.ID] = true
				}

				vars.Users = users
				vars.Members = members
			}
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
