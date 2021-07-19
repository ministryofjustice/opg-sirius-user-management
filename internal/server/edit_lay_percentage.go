package server

import (
	"fmt"
	"strconv"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type EditLayPercentageClient interface {
	EditLayPercentage(ctx sirius.Context, layPercentage string) (error)
	RandomReviews(sirius.Context) (sirius.RandomReviews, error)
}

type editLayPercentageVars struct {
	Path                string
	XSRFToken           string
	LayPercentage       string
	ReviewCycle         string
	Success             bool
	Errors              sirius.ValidationErrors
}

func editLayPercentage(client EditLayPercentageClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
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

            vars := editLayPercentageVars{
                Path:               r.URL.Path,
                XSRFToken:          ctx.XSRFToken,
                LayPercentage:      strconv.Itoa(randomReviews.LayPercentage),
            }

            return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
					fmt.Print("MethodPost")
            layPercentage := r.PostFormValue("layPercentage")

			err := client.EditLayPercentage(ctx, layPercentage)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars := editLayPercentageVars{
					Path:      r.URL.Path,
					XSRFToken: ctx.XSRFToken,
					LayPercentage: layPercentage,
					Errors:    verr.Errors,
				}

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			fmt.Print("hehre")
			fmt.Print(err)

			return RedirectError(fmt.Sprintf("/random-reviews/edit/lay-percentage"))

		default:
		fmt.Print("StatusError")
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
