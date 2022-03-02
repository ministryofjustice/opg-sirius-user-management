package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/template"
)

type ViewTeamClient interface {
	Team(sirius.Context, int) (sirius.Team, error)
}

type viewTeamVars struct {
	Path      string
	XSRFToken string
	Team      sirius.Team
}

func viewTeam(client ViewTeamClient, tmpl template.Template) handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		perm := myPermissions(r)
		if !perm.HasPermission("v1-teams", http.MethodPut) {
			return handler.Status(http.StatusForbidden)
		}

		if r.Method != http.MethodGet {
			return handler.Status(http.StatusMethodNotAllowed)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/teams/"))
		if err != nil {
			return handler.Status(http.StatusNotFound)
		}

		ctx := getContext(r)

		team, err := client.Team(ctx, id)
		if err != nil {
			return err
		}

		vars := viewTeamVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
			Team:      team,
		}

		return tmpl(w, vars)
	}
}
