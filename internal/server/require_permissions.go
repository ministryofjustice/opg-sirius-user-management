package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type RequirePermissionClient interface {
	GetMyPermissions(sirius.Context) (sirius.PermissionSet, error)
}

type PermissionRequest struct {
	group  string
	method string
}

func requirePermissions(client RequirePermissionClient, permissions ...PermissionRequest) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			myPermissions, err := client.GetMyPermissions(getContext(r))
			if err != nil {
				return err
			}

			for _, permission := range permissions {
				if !myPermissions.HasPermission(permission.group, permission.method) {
					return StatusError(http.StatusForbidden)
				}
			}

			return next(w, r)
		}
	}
}
