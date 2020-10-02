package server

import (
	"context"
	"log"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type ListUsersClient interface {
	SearchUsers(context.Context, []*http.Cookie, string) ([]sirius.User, error)
}

type listUsersVars struct {
	Path      string
	SiriusURL string

	Users  []sirius.User
	Search string
	Errors sirius.ValidationErrors
}

func listUsers(logger *log.Logger, client ListUsersClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		search := r.FormValue("search")

		vars := listUsersVars{
			Path:      r.URL.Path,
			SiriusURL: siriusURL,
			Search:    search,
		}

		if len(search) >= 3 {
			users, err := client.SearchUsers(r.Context(), r.Cookies(), search)
			if err != nil {
				return err
			}
			vars.Users = users
		} else if search != "" {
			vars.Errors = sirius.ValidationErrors{
				"search": {
					"": "Search term must be at least three characters",
				},
			}
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
