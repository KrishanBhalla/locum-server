package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/KrishanBhalla/locum-server/api/spec"
	"github.com/KrishanBhalla/locum-server/services"
	chiMw "github.com/go-chi/chi/middleware"
)

func FindUsers(ctx context.Context, request spec.FindUsersRequestObject) (spec.FindUsersResponseObject, error) {
	/// setup
	services, ok := services.FromContext(ctx)
	reqId := chiMw.GetReqID(ctx)
	internalServerError := spec.FindUsersdefaultResponse{StatusCode: http.StatusInternalServerError}
	if !ok {
		return internalServerError, errors.New(fmt.Sprintf("No services passed via context, reqId: %s", reqId))
	}

	// Process

	foundUsers, err := services.User.Query(request.Body.QueryString)
	if err != nil {
		return internalServerError, err
	}
	// Create Response
	resp := make(spec.FindUsers200JSONResponse, len(foundUsers), len(foundUsers))
	for i, u := range foundUsers {
		resp[i] = spec.UserResponse{UserId: u.Id, FullName: u.FullName}
	}
	return resp, nil
}
