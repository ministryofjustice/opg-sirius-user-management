package server

import (
	"github.com/ministryofjustice/opg-sirius-user-management/internal/model"
	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"net/http"
)

type FeedbackFormClient interface {
	AddFeedback(sirius.Context, model.FeedbackForm) error
}

type feedbackFormVars struct {
	Path    string
	Success bool
	Error   sirius.ValidationError
	Form    model.FeedbackForm
}

func feedbackForm(client FeedbackFormClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}
		ctx := getContext(r)
		vars := feedbackFormVars{
			Path: "/supervision/feedback",
		}

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			//if err != nil {
			//	return err
			//}

			err = client.AddFeedback(ctx, model.FeedbackForm{
				IsSupervisionFeedback: true,
				Name:                  r.FormValue("name"),
				Email:                 r.FormValue("email"),
				CaseNumber:            r.FormValue("case-number"),
				Message:               r.FormValue("more-detail"),
			})

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Error.Message = verr.Message
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			} else {
				vars.Success = true
			}
		}
		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
