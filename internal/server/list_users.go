package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/template"
)

type ListUsersClient interface {
	SearchUsers(sirius.Context, string) ([]sirius.User, error)
}

type listUsersVars struct {
	Path   string
	Users  []sirius.User
	Search string
	Errors sirius.ValidationErrors
}

func listUsers(client ListUsersClient, tmpl template.Template) handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		perm := myPermissions(r)

		if !perm.HasPermission("v1-users", http.MethodPut) {
			return handler.Status(http.StatusForbidden)
		}

		if r.Method != http.MethodGet {
			return handler.Status(http.StatusMethodNotAllowed)
		}

		search := r.FormValue("search")

		vars := listUsersVars{
			Path:   r.URL.Path,
			Search: search,
		}

		if search != "" {
			users, err := client.SearchUsers(getContext(r), search)

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

		return tmpl(w, vars)
	}
}
