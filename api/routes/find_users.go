package routes

import (
	"context"
	"net/http"

	"github.com/KrishanBhalla/locum-server/api/spec"
)

func FindUsers(ctx context.Context, request spec.FindUsersRequestObject) (spec.FindUsersResponseObject, error) {
	/// setup
	internalServerError := spec.FindUsersdefaultResponse{StatusCode: http.StatusInternalServerError}

	services, err := validateServices(ctx)
	if err != nil {
		return internalServerError, err
	}

	userToken, err := validateToken(ctx, services)
	if err != nil {
		return spec.UnauthorizedErrorResponse{}, err
	}

	// Process
	foundUsers, err := services.User.Query(request.Body.QueryString)
	if err != nil {
		return internalServerError, err
	}
	// Create Response
	resp := make(spec.FindUsers200JSONResponse, 0, len(foundUsers))
	for _, u := range foundUsers {
		if u.Id != userToken.UserId {
			resp = append(resp, spec.UserResponse{UserId: u.Id, FullName: u.FullName})
		}
	}
	return resp, nil
}
