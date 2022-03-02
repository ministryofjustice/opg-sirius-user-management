package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/template"
)

type EditTeamClient interface {
	Team(sirius.Context, int) (sirius.Team, error)
	EditTeam(sirius.Context, sirius.Team) error
	TeamTypes(sirius.Context) ([]sirius.RefDataTeamType, error)
}

type editTeamVars struct {
	Path            string
	XSRFToken       string
	Team            sirius.Team
	TeamTypeOptions []sirius.RefDataTeamType
	CanEditTeamType bool
	CanDeleteTeam   bool
	Success         bool
	Errors          sirius.ValidationErrors
}

func editTeam(client EditTeamClient, tmpl template.Template) handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		perm := myPermissions(r)
		if !perm.HasPermission("v1-teams", http.MethodPut) {
			return handler.Status(http.StatusForbidden)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/teams/edit/"))
		if err != nil {
			return handler.Status(http.StatusNotFound)
		}

		ctx := getContext(r)

		team, err := client.Team(ctx, id)
		if err != nil {
			return err
		}

		canEditTeamType := perm.HasPermission("v1-teams", http.MethodPost)
		canDeleteTeam := perm.HasPermission("v1-teams", http.MethodDelete)

		teamTypes, err := client.TeamTypes(ctx)
		if err != nil {
			return err
		}

		vars := editTeamVars{
			Path:            r.URL.Path,
			XSRFToken:       ctx.XSRFToken,
			Team:            team,
			TeamTypeOptions: teamTypes,
			CanEditTeamType: canEditTeamType,
			CanDeleteTeam:   canDeleteTeam,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl(w, vars)
		case http.MethodPost:
			vars.Team.DisplayName = r.PostFormValue("name")
			vars.Team.PhoneNumber = r.PostFormValue("phone")
			vars.Team.Email = r.PostFormValue("email")

			if canEditTeamType {
				if r.PostFormValue("service") == "supervision" {
					vars.Team.Type = r.PostFormValue("supervision-type")
				} else {
					vars.Team.Type = ""
				}
			} else {
				vars.Team.Type = team.Type
			}

			// Attempt to save
			err := client.EditTeam(ctx, vars.Team)

			if e, ok := err.(*sirius.ValidationError); ok {
				vars.Errors = e.Errors
				w.WriteHeader(http.StatusBadRequest)
			} else if err != nil {
				return err
			} else {
				vars.Success = true
			}

			return tmpl(w, vars)
		default:
			return handler.Status(http.StatusMethodNotAllowed)
		}
	}
}
