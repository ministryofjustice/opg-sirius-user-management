package server

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type ListUsersClient interface {
	ListUsers(context.Context, []*http.Cookie) ([]sirius.User, error)
}

type listUsersVars struct {
	Path      string
	SiriusURL string

	Users    []sirius.User
	Search   string
	PagePrev int64
	PageNext int64
}

func prepareSearchTerm(term string) string {
	return strings.Replace(strings.ToLower(term), " ", "", -1)
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

		pageSize := int64(50)
		search := r.FormValue("search")
		page, err := strconv.ParseInt(r.FormValue("page"), 10, 64)

		if err != nil {
			page = 1
		}

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

		pageStart := (page - 1) * pageSize
		pageEnd := pageStart + pageSize

		if pageStart > int64(len(filtered)) {
			return StatusError(http.StatusNotFound)
		}

		if pageEnd > int64(len(filtered)) {
			pageEnd = int64(len(filtered))
		}

		var pagePrev int64
		var pageNext int64

		if page > 1 {
			pagePrev = page - 1
		} else {
			pagePrev = 0
		}

		if page*pageSize > int64(len(filtered)) {
			pageNext = 0
		} else {
			pageNext = page + 1
		}

		vars := listUsersVars{
			Path:      r.URL.Path,
			SiriusURL: siriusURL,
			Users:     filtered[pageStart:pageEnd],
			Search:    search,
			PagePrev:  pagePrev,
			PageNext:  pageNext,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
