package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/template"
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

func editLayPercentage(client EditLayPercentageClient, tmpl template.Template) handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		perm := myPermissions(r)
		if !perm.HasPermission("v1-random-review-settings", http.MethodPost) {
			return handler.Status(http.StatusForbidden)
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

			return tmpl(w, vars)

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
				return tmpl(w, vars)
			} else if err != nil {
				return err
			}

			return handler.Redirect("/random-reviews")

		default:
			return handler.Status(http.StatusMethodNotAllowed)
		}
	}
}
