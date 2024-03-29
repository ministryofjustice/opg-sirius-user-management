package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type EditLayPercentageClient interface {
	RandomReviews(sirius.Context) (sirius.RandomReviews, error)
	EditRandomReviewSettings(ctx sirius.Context, layPercentage string, reviewCycle string) error
}

type editLayPercentageVars struct {
	Path          string
	XSRFToken     string
	LayPercentage string
	ReviewCycle   string
	Errors        sirius.ValidationErrors
	Error         string
}

func editLayPercentage(client EditLayPercentageClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if !perm.HasPermission("v1-random-review-settings", http.MethodPost) {
			return StatusError(http.StatusForbidden)
		}

		ctx := getContext(r)

		vars := editLayPercentageVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
		}

		switch r.Method {
		case http.MethodGet:
			randomReviews, err := client.RandomReviews(ctx)
			if err != nil {
				return err
			}

			vars.LayPercentage = strconv.Itoa(randomReviews.LayPercentage)

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:

			randomReviewsCycle, _ := client.RandomReviews(ctx)

			layPercentage := r.PostFormValue("layPercentage")
			reviewCycle := strconv.Itoa(randomReviewsCycle.ReviewCycle)

			err := client.EditRandomReviewSettings(ctx, reviewCycle, layPercentage)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.LayPercentage = layPercentage
				vars.Errors = verr.Errors
				vars.Error = verr.Message
				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return RedirectError("/random-reviews")

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
