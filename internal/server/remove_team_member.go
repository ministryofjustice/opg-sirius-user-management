package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type RemoveTeamMemberClient interface {
	Team(context.Context, []*http.Cookie, int) (sirius.Team, error)
	EditTeam(context.Context, []*http.Cookie, sirius.Team) error
}

type removeTeamMemberVars struct {
	Path      string
	SiriusURL string
	Team      sirius.Team
	Selected  []sirius.TeamMember
	Errors    sirius.ValidationErrors
}

func removeTeamMember(client RemoveTeamMemberClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/teams/remove-member/"))
		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		team, err := client.Team(r.Context(), r.Cookies(), id)
		if err != nil {
			return err
		}

		vars := removeTeamMemberVars{
			Path:      r.URL.Path,
			SiriusURL: siriusURL,
			Team:      team,
		}

		r.ParseForm()
		for _, id := range r.PostForm["selected[]"] {
			userID, err := strconv.Atoi(id)
			if err != nil {
				return StatusError(http.StatusBadRequest)
			}

			for _, user := range team.Members {
				if userID == user.ID {
					vars.Selected = append(vars.Selected, user)
				}
			}
		}

		if r.PostFormValue("confirm") != "" {
			var members []sirius.TeamMember
			for _, member := range team.Members {
				stillInTeam := true
				for _, selected := range vars.Selected {
					if member.ID == selected.ID {
						stillInTeam = false
					}
				}

				if stillInTeam {
					members = append(members, member)
				}
			}

			team.Members = members

			err = client.EditTeam(r.Context(), r.Cookies(), team)

			if _, ok := err.(sirius.ClientError); ok {
				vars.Errors = sirius.ValidationErrors{
					"_": {
						"": err.Error(),
					},
				}
				w.WriteHeader(http.StatusBadRequest)
			} else if err != nil {
				return err
			} else {
				return RedirectError(fmt.Sprintf("/teams/%d", team.ID))
			}
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
