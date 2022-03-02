package server

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/securityheaders"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/template"
)

type Logger interface {
	Request(*http.Request, error)
}

type Client interface {
	AddTeamClient
	AddUserClient
	ChangePasswordClient
	DeleteTeamClient
	DeleteUserClient
	EditMyDetailsClient
	EditTeamClient
	EditUserClient
	MyPermissionsClient
	ListTeamsClient
	ListUsersClient
	MyDetailsClient
	ResendConfirmationClient
	UnlockUserClient
	ViewTeamClient
	RandomReviewsClient
	EditRandomReviewSettingsClient
}

func New(logger Logger, client Client, templates template.Templates, prefix, siriusPublicURL, webDir string) http.Handler {
	errorTmpl := templates.Get("error.gotmpl")

	wrap := handler.New(prefix, siriusPublicURL+"/auth", func(w http.ResponseWriter, r *http.Request, code int, err error) {
		logger.Request(r, err)

		err = errorTmpl(w, errorVars{
			SiriusURL: siriusPublicURL,
			Path:      "",
			Code:      code,
			Error:     err.Error(),
		})

		if err != nil {
			logger.Request(r, err)
			http.Error(w, "Could not generate error template", http.StatusInternalServerError)
		}
	}, withMyPermissions(client))

	mux := http.NewServeMux()
	mux.Handle("/", http.RedirectHandler(prefix+"/my-details", http.StatusFound))
	mux.Handle("/health-check", healthCheck())

	mux.Handle("/users",
		wrap(
			listUsers(client, templates.Get("users.gotmpl"))))

	mux.Handle("/teams",
		wrap(
			listTeams(client, templates.Get("teams.gotmpl"))))

	mux.Handle("/teams/",
		wrap(
			viewTeam(client, templates.Get("team.gotmpl"))))

	mux.Handle("/teams/add",
		wrap(
			addTeam(client, templates.Get("team-add.gotmpl"))))

	mux.Handle("/teams/edit/",
		wrap(
			editTeam(client, templates.Get("team-edit.gotmpl"))))

	mux.Handle("/teams/delete/",
		wrap(
			deleteTeam(client, templates.Get("team-delete.gotmpl"))))

	mux.Handle("/teams/add-member/",
		wrap(
			addTeamMember(client, templates.Get("team-add-member.gotmpl"))))

	mux.Handle("/teams/remove-member/",
		wrap(
			removeTeamMember(client, templates.Get("team-remove-member.gotmpl"))))

	mux.Handle("/my-details",
		wrap(
			myDetails(client, templates.Get("my-details.gotmpl"))))

	mux.Handle("/my-details/edit",
		wrap(
			editMyDetails(client, templates.Get("edit-my-details.gotmpl"))))

	mux.Handle("/random-reviews",
		wrap(
			randomReviews(client, templates.Get("random-reviews.gotmpl"))))

	mux.Handle("/random-reviews/edit/lay-percentage",
		wrap(
			editRandomReviewSettings(client, templates.Get("random-reviews-edit-lay-percentage.gotmpl"))))

	mux.Handle("/random-reviews/edit/pa-percentage",
		wrap(
			editRandomReviewSettings(client, templates.Get("random-reviews-edit-pa-percentage.gotmpl"))))

	mux.Handle("/random-reviews/edit/pro-percentage",
		wrap(
			editRandomReviewSettings(client, templates.Get("random-reviews-edit-pro-percentage.gotmpl"))))

	mux.Handle("/random-reviews/edit/review-cycle",
		wrap(
			editRandomReviewSettings(client, templates.Get("random-reviews-edit-review-cycle.gotmpl"))))

	mux.Handle("/change-password",
		wrap(
			changePassword(client, templates.Get("change-password.gotmpl"))))

	mux.Handle("/add-user",
		wrap(
			addUser(client, templates.Get("add-user.gotmpl"))))

	mux.Handle("/edit-user/",
		wrap(
			editUser(client, templates.Get("edit-user.gotmpl"))))

	mux.Handle("/unlock-user/",
		wrap(
			unlockUser(client, templates.Get("unlock-user.gotmpl"))))

	mux.Handle("/delete-user/",
		wrap(
			deleteUser(client, templates.Get("delete-user.gotmpl"))))

	mux.Handle("/resend-confirmation",
		wrap(
			resendConfirmation(client, templates.Get("resend-confirmation.gotmpl"))))

	static := http.FileServer(http.Dir(webDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	return http.StripPrefix(prefix, securityheaders.Use(mux))
}

type errorVars struct {
	SiriusURL string
	Path      string

	Code  int
	Error string
}

type MyPermissionsClient interface {
	MyPermissions(sirius.Context) (sirius.PermissionSet, error)
}

type myPermissionsKey struct{}

func withMyPermissions(client MyPermissionsClient) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			myPermissions, err := client.MyPermissions(getContext(r))
			if err != nil {
				return err
			}

			return next(w, r.WithContext(context.WithValue(r.Context(), myPermissionsKey{}, myPermissions)))
		}
	}
}

func myPermissions(r *http.Request) sirius.PermissionSet {
	return r.Context().Value(myPermissionsKey{}).(sirius.PermissionSet)
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
