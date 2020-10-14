package server

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type ViewTeamClient interface {
	Team(context.Context, []*http.Cookie, int) (sirius.Team, error)
}

type viewTeamVars struct {
	Path      string
	SiriusURL string
	Team      sirius.Team
}

func viewTeam(client ViewTeamClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/teams/"))
		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		team, err := client.Team(r.Context(), r.Cookies(), id)
		if err != nil {
			return err
		}

		vars := viewTeamVars{
			Path:      r.URL.Path,
			SiriusURL: siriusURL,
			Team:      team,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
