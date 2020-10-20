package server

import (
	"context"
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

func listUsers(client ListUsersClient, tmpl Template, siriusURL string) Handler {
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

		if search != "" {
			users, err := client.SearchUsers(r.Context(), r.Cookies(), search)

			if _, ok := err.(sirius.ClientError); ok {
				vars.Errors = sirius.ValidationErrors{
					"search": {
						"": err.Error(),
					},
				}
			} else if err != nil {
				return err
			} else {
				vars.Users = users
			}
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
