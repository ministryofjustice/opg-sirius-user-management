package server

import (
	"github.com/ministryofjustice/opg-sirius-user-management/internal/model"
	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"net/http"
)

type FeedbackFormClient interface {
	SubmitFeedback(sirius.Context, model.FeedbackForm) error
	MyDetails(sirius.Context) (sirius.MyDetails, error)
}

type feedbackFormVars struct {
	ID             int
	SuccessMessage string
	Error          sirius.ValidationError
	Form           model.FeedbackForm
}

func feedbackForm(client FeedbackFormClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}
		ctx := getContext(r)
		vars := feedbackFormVars{}

		if r.Method == http.MethodGet {

			return tmpl.ExecuteTemplate(w, "page", vars)
		}

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				return err
			}

			if len(r.FormValue("more-detail")) < 1 {
				vars.Error = sirius.ValidationError{
					Message: "no-feedback",
					Errors:  nil,
				}
			}
			if len(r.FormValue("more-detail")) > 900 {
				vars.Error = sirius.ValidationError{
					Message: "feedback-too-long",
					Errors:  nil,
				}
			}

			if vars.Error.Message != "" {
				vars.Form = model.FeedbackForm{
					IsSupervisionFeedback: true,
					Name:                  r.FormValue("name"),
					Email:                 r.FormValue("email"),
					CaseNumber:            r.FormValue("case-number"),
					Message:               r.FormValue("more-detail"),
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			err = client.SubmitFeedback(ctx, model.FeedbackForm{
				IsSupervisionFeedback: true,
				Name:                  r.FormValue("name"),
				Email:                 r.FormValue("email"),
				CaseNumber:            r.FormValue("case-number"),
				Message:               r.FormValue("more-detail"),
			})

			if err != nil {
				return err
			}

			vars.SuccessMessage = "Form Submitted"
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
