package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/template"
)

type RemoveTeamMemberClient interface {
	Team(sirius.Context, int) (sirius.Team, error)
	EditTeam(sirius.Context, sirius.Team) error
}

type removeTeamMemberVars struct {
	Path      string
	XSRFToken string
	Team      sirius.Team
	Selected  map[int]string
	Errors    sirius.ValidationErrors
}

func removeTeamMember(client RemoveTeamMemberClient, tmpl template.Template) handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		perm := myPermissions(r)
		if !perm.HasPermission("v1-teams", http.MethodPut) {
			return handler.Status(http.StatusForbidden)
		}

		if r.Method != http.MethodPost {
			return handler.Status(http.StatusMethodNotAllowed)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/teams/remove-member/"))
		if err != nil {
			return handler.Status(http.StatusNotFound)
		}

		if err := r.ParseForm(); err != nil {
			return handler.Status(http.StatusBadRequest)
		}

		ctx := getContext(r)

		team, err := client.Team(ctx, id)
		if err != nil {
			return err
		}

		vars := removeTeamMemberVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
			Team:      team,
			Selected:  make(map[int]string),
		}

		for _, id := range r.PostForm["selected[]"] {
			userID, err := strconv.Atoi(id)
			if err != nil {
				return handler.Status(http.StatusBadRequest)
			}

			for _, user := range team.Members {
				if userID == user.ID {
					vars.Selected[userID] = user.DisplayName
				}
			}
		}

		if r.PostFormValue("confirm") != "" {
			var members []sirius.TeamMember
			for _, member := range team.Members {
				if _, ok := vars.Selected[member.ID]; !ok {
					members = append(members, member)
				}
			}

			team.Members = members

			err = client.EditTeam(ctx, team)

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
				return handler.Redirect(fmt.Sprintf("/teams/%d", team.ID))
			}
		}

		return tmpl(w, vars)
	}
}
