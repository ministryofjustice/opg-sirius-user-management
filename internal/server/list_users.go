package server

import (
	"context"
	"log"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type ListUsersClient interface {
	ListUsers(context.Context, []*http.Cookie) ([]sirius.User, error)
}

type listUsersVars struct {
	Path      string
	SiriusURL string

	Users []sirius.User
}

func listUsers(logger *log.Logger, client ListUsersClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		users, err := client.ListUsers(r.Context(), r.Cookies())
		if err != nil {
			return err
		}

		vars := listUsersVars{
			Path:      r.URL.Path,
			SiriusURL: siriusURL,
			Users:     users,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
