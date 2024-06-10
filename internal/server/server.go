package server

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Logger interface {
	Request(*http.Request, error)
}

type Client interface {
	AddTeamClient
	AddUserClient
	DeleteTeamClient
	DeleteUserClient
	EditMyDetailsClient
	EditTeamClient
	EditUserClient
	ErrorHandlerClient
	ListTeamsClient
	ListUsersClient
	MyDetailsClient
	ViewTeamClient
	RandomReviewsClient
	EditRandomReviewSettingsClient
	FeedbackFormClient
}

type Template interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(logger *slog.Logger, client Client, templates map[string]*template.Template, prefix, siriusPublicURL, webDir string) http.Handler {
	wrap := errorHandler(client, templates["error.gotmpl"], prefix, siriusPublicURL)

	mux := http.NewServeMux()
	mux.Handle("/", http.RedirectHandler(prefix+"/my-details", http.StatusFound))
	mux.Handle("/health-check", healthCheck())

	mux.Handle("/users",
		wrap(
			listUsers(client, templates["users.gotmpl"])))

	mux.Handle("/teams",
		wrap(
			listTeams(client, templates["teams.gotmpl"])))

	mux.Handle("/teams/",
		wrap(
			viewTeam(client, templates["team.gotmpl"])))

	mux.Handle("/teams/add",
		wrap(
			addTeam(client, templates["team-add.gotmpl"])))

	mux.Handle("/teams/edit/",
		wrap(
			editTeam(client, templates["team-edit.gotmpl"])))

	mux.Handle("/teams/delete/",
		wrap(
			deleteTeam(client, templates["team-delete.gotmpl"])))

	mux.Handle("/teams/add-member/",
		wrap(
			addTeamMember(client, templates["team-add-member.gotmpl"])))

	mux.Handle("/teams/remove-member/",
		wrap(
			removeTeamMember(client, templates["team-remove-member.gotmpl"])))

	mux.Handle("/my-details",
		wrap(
			myDetails(client, templates["my-details.gotmpl"])))

	mux.Handle("/my-details/edit",
		wrap(
			editMyDetails(client, templates["edit-my-details.gotmpl"])))

	mux.Handle("/random-reviews",
		wrap(
			randomReviews(client, templates["random-reviews.gotmpl"])))

	mux.Handle("/random-reviews/edit/lay-percentage",
		wrap(
			editRandomReviewSettings(client, templates["random-reviews-edit-lay-percentage.gotmpl"])))

	mux.Handle("/random-reviews/edit/pa-percentage",
		wrap(
			editRandomReviewSettings(client, templates["random-reviews-edit-pa-percentage.gotmpl"])))

	mux.Handle("/random-reviews/edit/pro-percentage",
		wrap(
			editRandomReviewSettings(client, templates["random-reviews-edit-pro-percentage.gotmpl"])))

	mux.Handle("/random-reviews/edit/review-cycle",
		wrap(
			editRandomReviewSettings(client, templates["random-reviews-edit-review-cycle.gotmpl"])))

	mux.Handle("/add-user",
		wrap(
			addUser(client, templates["add-user.gotmpl"])))

	mux.Handle("/edit-user/",
		wrap(
			editUser(client, templates["edit-user.gotmpl"])))

	mux.Handle("/delete-user/",
		wrap(
			deleteUser(client, templates["delete-user.gotmpl"])))

	mux.Handle("/supervision/feedback",
		wrap(
			feedbackForm(client, templates["feedback.gotmpl"])))

	static := http.FileServer(http.Dir(webDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	middleware := telemetry.Middleware(logger)

	return otelhttp.NewHandler(http.StripPrefix(prefix, securityheaders.Use(middleware(mux))), "user-management")
}

type RedirectError string

func (e RedirectError) Error() string {
	return "redirect to " + string(e)
}

func (e RedirectError) To() string {
	return string(e)
}

type StatusError int

func (e StatusError) Error() string {
	code := e.Code()

	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}

func (e StatusError) Code() int {
	return int(e)
}

type Handler func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error

type errorVars struct {
	SiriusURL string
	Path      string

	Code  int
	Error string
}

type ErrorHandlerClient interface {
	MyPermissions(sirius.Context) (sirius.PermissionSet, error)
}

func errorHandler(client ErrorHandlerClient, tmplError Template, prefix, siriusURL string) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			myPermissions, err := client.MyPermissions(getContext(r))

			if err == nil {
				err = next(myPermissions, w, r)
			}

			if err != nil {
				if errors.Is(err, context.Canceled) {
					w.WriteHeader(499)
					return
				}

				if err == sirius.ErrUnauthorized {
					http.Redirect(w, r, fmt.Sprintf("%s/auth?redirect=%s", siriusURL, url.QueryEscape(prefix+r.URL.Path)), http.StatusFound)
					return
				}

				if redirect, ok := err.(RedirectError); ok {
					http.Redirect(w, r, prefix+redirect.To(), http.StatusFound)
					return
				}

				code := http.StatusInternalServerError
				if status, ok := err.(StatusError); ok {
					if status.Code() == http.StatusForbidden || status.Code() == http.StatusNotFound {
						code = status.Code()
					}
				}

				logger := telemetry.LoggerFromContext(r.Context())
				if code == http.StatusInternalServerError {
					logger.Error(err.Error())
				}

				w.WriteHeader(code)
				err = tmplError.ExecuteTemplate(w, "page", errorVars{
					SiriusURL: siriusURL,
					Path:      "",
					Code:      code,
					Error:     err.Error(),
				})

				if err != nil {
					logger.Error("could not generate error template", slog.Any("err", err.Error()))
					http.Error(w, "Could not generate error template", http.StatusInternalServerError)
				}
			}
		})
	}
}

func getContext(r *http.Request) sirius.Context {
	token := ""

	if r.Method == http.MethodGet {
		if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
			token, _ = url.QueryUnescape(cookie.Value)
		}
	} else {
		token = r.FormValue("xsrfToken")
	}

	return sirius.Context{
		Context:   r.Context(),
		Cookies:   r.Cookies(),
		XSRFToken: token,
	}
}
