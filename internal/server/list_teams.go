package server

import (
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/template"
)

type ListTeamsClient interface {
	Teams(sirius.Context) ([]sirius.Team, error)
}

type listTeamsVars struct {
	Path   string
	Search string
	Teams  []sirius.Team
}

func listTeams(client ListTeamsClient, tmpl template.Template) handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		perm := myPermissions(r)
		if !perm.HasPermission("v1-teams", http.MethodPut) {
			return handler.Status(http.StatusForbidden)
		}

		if r.Method != http.MethodGet {
			return handler.Status(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		teams, err := client.Teams(ctx)
		if err != nil {
			return err
		}

		search := r.FormValue("search")
		if search != "" {
			searchLower := strings.ToLower(search)

			var matchingTeams []sirius.Team
			for _, t := range teams {
				if strings.Contains(strings.ToLower(t.DisplayName), searchLower) {
					matchingTeams = append(matchingTeams, t)
				}
			}

			teams = matchingTeams
		}

		vars := listTeamsVars{
			Path:   r.URL.Path,
			Search: search,
			Teams:  teams,
		}

		return tmpl(w, vars)
	}
}
