package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type EditLayReviewCycleClient interface {
	EditLayReviewCycle(ctx sirius.Context, reviewCycle string, layPercentage string) (error)
	RandomReviews(sirius.Context) (sirius.RandomReviews, error)
}

type editLayReviewCycleVars struct {
	Path                string
	XSRFToken           string
	LayPercentage       string
	ReviewCycle         string
	Success             bool
	Errors              sirius.ValidationErrors
	Error               string
}

func editLayReviewCycle(client EditLayReviewCycleClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error{
		if !perm.HasPermission("v1-random-review-settings", http.MethodPost) {
			return StatusError(http.StatusForbidden)
		}

		ctx := getContext(r)

		switch r.Method {
		case http.MethodGet:
            randomReviews, err := client.RandomReviews(ctx)
            if err != nil {
                return err
            }

            vars := editLayReviewCycleVars{
                Path:               r.URL.Path,
                XSRFToken:          ctx.XSRFToken,
                ReviewCycle:        strconv.Itoa(randomReviews.ReviewCycle),
            }

            return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
		    randomReviewsCycle, _ := client.RandomReviews(ctx)

            reviewCycle := r.PostFormValue("layReviewCycle")
            layPercentage := strconv.Itoa(randomReviewsCycle.LayPercentage)

			err := client.EditLayReviewCycle(ctx, reviewCycle, layPercentage)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars := editLayReviewCycleVars{
					ReviewCycle:    reviewCycle,
					Errors:         verr.Errors,
					Error:          verr.Message,
				}

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}
			return Redirect(fmt.Sprintf("/random-reviews"))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
