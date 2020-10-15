package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type ListTeamsClient interface {
	Teams(context.Context, []*http.Cookie) ([]sirius.Team, error)
}

type listTeamsVars struct {
	Path      string
	SiriusURL string
	Search    string
	Teams     []sirius.Team
}

func listTeams(client ListTeamsClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		teams, err := client.Teams(r.Context(), r.Cookies())
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
			Path:      r.URL.Path,
			SiriusURL: siriusURL,
			Search:    search,
			Teams:     teams,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
