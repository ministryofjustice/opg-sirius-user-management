package server

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type ListUsersClient interface {
	ListUsers(context.Context, []*http.Cookie) ([]sirius.User, error)
	MyDetails(context.Context, []*http.Cookie) (sirius.MyDetails, error)
}

type listUsersVars struct {
	Path      string
	SiriusURL string

	Users  []sirius.User
	Search string
}

func prepareSearchTerm(term string) string {
	return strings.Replace(strings.ToLower(term), " ", "", -1)
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

		users, err := client.ListUsers(r.Context(), r.Cookies())
		if err != nil {
			return err
		}

		search := r.FormValue("search")

		var filtered []sirius.User

		if search != "" {
			preparedSearch := prepareSearchTerm(search)

			for _, user := range users {
				if strings.Contains(prepareSearchTerm(user.DisplayName), preparedSearch) ||
					strings.Contains(prepareSearchTerm(user.Email), preparedSearch) ||
					strings.Contains(prepareSearchTerm(user.Status.String()), preparedSearch) {
					filtered = append(filtered, user)
				}
			}
		} else {
			filtered = users
		}

		vars := listUsersVars{
			Path:      r.URL.Path,
			SiriusURL: siriusURL,
			Users:     filtered,
			Search:    search,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
