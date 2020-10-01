package server

import (
	"context"
	"log"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type ListUsersClient interface {
	SearchUsers(context.Context, []*http.Cookie, string) ([]sirius.User, error)
	MyDetails(context.Context, []*http.Cookie) (sirius.MyDetails, error)
}

type listUsersVars struct {
	Path      string
	SiriusURL string

	Users  []sirius.User
	Search string
	Error  string
}

func listUsers(logger *log.Logger, client ListUsersClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		myDetails, err := client.MyDetails(r.Context(), r.Cookies())
		if err != nil {
			return err
		}

		permitted := false
		for _, role := range myDetails.Roles {
			if role == "System Admin" {
				permitted = true
			}
		}

		if !permitted {
			return StatusError(http.StatusForbidden)
		}

		search := r.FormValue("search")
		var users []sirius.User

		errorMessage := ""

		if search != "" && len(search) < 3 {
			errorMessage = "Search term must be at least three characters"
		} else if search != "" {
			var err error
			users, err = client.SearchUsers(r.Context(), r.Cookies(), search)
			if err != nil {
				return err
			}
		}

		vars := listUsersVars{
			Path:      r.URL.Path,
			SiriusURL: siriusURL,
			Users:     users,
			Search:    search,
			Error:     errorMessage,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
