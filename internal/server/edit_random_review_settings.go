package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type EditRandomReviewSettingsClient interface {
	EditRandomReviewSettings(ctx sirius.Context, randomReviews sirius.EditRandomReview) error
	RandomReviews(sirius.Context) (sirius.RandomReviews, error)
}

type editRandomReviewSettingsVars struct {
	Path          string
	XSRFToken     string
	LayPercentage int
	PaPercentage  int
	ReviewCycle   int
	Errors        sirius.ValidationErrors
	Error         string
}

func editRandomReviewSettings(client EditRandomReviewSettingsClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if !perm.HasPermission("v1-random-review-settings", http.MethodPost) {
			return StatusError(http.StatusForbidden)
		}

		ctx := getContext(r)

		vars := editRandomReviewSettingsVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
		}

		switch r.Method {
		case http.MethodGet:
			randomReviews, err := client.RandomReviews(ctx)
			if err != nil {
				return err
			}

			vars.LayPercentage = randomReviews.LayPercentage
			vars.PaPercentage = randomReviews.PaPercentage
			vars.ReviewCycle = randomReviews.ReviewCycle

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			randomReviewSettings, _ := client.RandomReviews(ctx)

			err := client.EditRandomReviewSettings(ctx, formValueOrExisting(r, randomReviewSettings))

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.LayPercentage = randomReviewSettings.LayPercentage
				vars.PaPercentage = randomReviewSettings.PaPercentage
				vars.ReviewCycle = randomReviewSettings.ReviewCycle
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

func formValueOrExisting(r *http.Request, existing sirius.RandomReviews) sirius.EditRandomReview {
	layPercentage := r.PostFormValue("layPercentage")
	if layPercentage == "" {
		layPercentage = strconv.Itoa(existing.LayPercentage)
	}

	paPercentage := r.PostFormValue("paPercentage")
	if paPercentage == "" {
		paPercentage = strconv.Itoa(existing.PaPercentage)
	}

	reviewCycle := r.PostFormValue("reviewCycle")
	if reviewCycle == "" {
		reviewCycle = strconv.Itoa(existing.ReviewCycle)
	}

	return sirius.EditRandomReview{
		LayPercentage: layPercentage,
		PaPercentage: paPercentage,
		ReviewCycle: reviewCycle,
	}
}
