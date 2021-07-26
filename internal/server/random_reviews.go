package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type RandomReviewsClient interface {
	RandomReviews(sirius.Context) (sirius.RandomReviews, error)
}

type randomReviewsVars struct {
	Path          string
	LayPercentage int
	ReviewCycle   int
}

func randomReviews(client RandomReviewsClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if !perm.HasPermission("v1-random-review-settings", http.MethodGet) {
			return StatusError(http.StatusForbidden)
		}

		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		randomReviews, err := client.RandomReviews(ctx)
		if err != nil {
			return err
		}

		vars := randomReviewsVars{
			Path:          r.URL.Path,
			LayPercentage: randomReviews.LayPercentage,
			ReviewCycle:   randomReviews.ReviewCycle,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
