package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/template"
)

type RandomReviewsClient interface {
	RandomReviews(sirius.Context) (sirius.RandomReviews, error)
}

type randomReviewsVars struct {
	Path          string
	LayPercentage int
	PaPercentage  int
	ProPercentage int
	ReviewCycle   int
}

func randomReviews(client RandomReviewsClient, tmpl template.Template) handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		perm := myPermissions(r)
		if !perm.HasPermission("v1-random-review-settings", http.MethodGet) {
			return handler.Status(http.StatusForbidden)
		}

		if r.Method != http.MethodGet {
			return handler.Status(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		randomReviews, err := client.RandomReviews(ctx)
		if err != nil {
			return err
		}

		vars := randomReviewsVars{
			Path:          r.URL.Path,
			LayPercentage: randomReviews.LayPercentage,
			PaPercentage:  randomReviews.PaPercentage,
			ProPercentage: randomReviews.ProPercentage,
			ReviewCycle:   randomReviews.ReviewCycle,
		}

		return tmpl(w, vars)
	}
}
