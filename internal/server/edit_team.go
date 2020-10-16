package server

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type EditTeamClient interface {
	Team(context.Context, []*http.Cookie, int) (sirius.Team, error)
	EditTeam(context.Context, []*http.Cookie, sirius.Team) error
	TeamTypes(context.Context, []*http.Cookie) ([]sirius.RefDataTeamType, error)
}

type editTeamVars struct {
	Path            string
	SiriusURL       string
	Team            sirius.Team
	TeamTypeOptions []sirius.RefDataTeamType
	Success         bool
	Errors          sirius.ValidationErrors
}

func editTeam(client EditTeamClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/teams/edit/"))
		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		team, err := client.Team(r.Context(), r.Cookies(), id)
		if err != nil {
			return err
		}

		teamTypes, err := client.TeamTypes(r.Context(), r.Cookies())
		if err != nil {
			return err
		}

		vars := editTeamVars{
			Path:            r.URL.Path,
			SiriusURL:       siriusURL,
			Team:            team,
			TeamTypeOptions: teamTypes,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)
		case http.MethodPost:
			vars.Team.DisplayName = r.PostFormValue("name")
			vars.Team.PhoneNumber = r.PostFormValue("phone")
			vars.Team.Email = r.PostFormValue("email")

			if r.PostFormValue("service") == "supervision" {
				vars.Team.Type = r.PostFormValue("supervision-type")
			} else {
				vars.Team.Type = ""
			}

			// Attempt to save
			err := client.EditTeam(r.Context(), r.Cookies(), vars.Team)

			if e, ok := err.(*sirius.ValidationError); ok {
				vars.Errors = e.Errors
				w.WriteHeader(http.StatusBadRequest)
			} else if err != nil {
				return err
			} else {
				vars.Success = true
			}

			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
