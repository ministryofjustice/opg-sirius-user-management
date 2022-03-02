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

type DeleteTeamClient interface {
	Team(sirius.Context, int) (sirius.Team, error)
	DeleteTeam(sirius.Context, int) error
}

type deleteTeamVars struct {
	Path           string
	XSRFToken      string
	Team           sirius.Team
	Errors         sirius.ValidationErrors
	SuccessMessage string
}

func deleteTeam(client DeleteTeamClient, tmpl template.Template) handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		perm := myPermissions(r)
		if !perm.HasPermission("v1-teams", http.MethodDelete) {
			return handler.Status(http.StatusForbidden)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/teams/delete/"))
		if err != nil {
			return handler.Status(http.StatusNotFound)
		}

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return handler.Status(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		team, err := client.Team(ctx, id)
		if err != nil {
			return err
		}

		vars := deleteTeamVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
			Team:      team,
		}

		if r.Method == http.MethodPost {
			err := client.DeleteTeam(ctx, id)

			if _, ok := err.(sirius.ClientError); ok {
				vars.Errors = sirius.ValidationErrors{
					"": {
						"": err.Error(),
					},
				}

				w.WriteHeader(http.StatusBadRequest)
			} else if err != nil {
				return err
			} else {
				vars.SuccessMessage = fmt.Sprintf("The team \"%s\" was deleted.", team.DisplayName)
			}
		}

		return tmpl(w, vars)
	}
}
