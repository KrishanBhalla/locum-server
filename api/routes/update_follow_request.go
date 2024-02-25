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

func UpdateFollowRequest(ctx context.Context, request spec.UpdateFollowRequestRequestObject) (spec.UpdateFollowRequestResponseObject, error) {

	services, ok := services.FromContext(ctx)
	reqId := chiMw.GetReqID(ctx)
	internalServerError := spec.UpdateFollowRequestdefaultResponse{StatusCode: http.StatusInternalServerError}
	if !ok {
		return internalServerError, errors.New(fmt.Sprintf("No services passed via context, reqId: %s", reqId))
	}

	// Process
	userId := request.Body.UserId
	requestingUser := request.Body.RequestedFollowerUserId
	requestAccepted := request.Body.Accept

	userFriends := services.UserFriends
	var err error
	if requestAccepted {
		err = userFriends.AddFollower(userId, requestingUser)
	}
	if err != nil {
		return nil, err
	}
	err = userFriends.RemoveFollowerRequest(userId, requestingUser)
	if err != nil {
		return nil, err
	}

	err = userFriends.RemoveFollowRequest(requestingUser, userId)
	if err != nil {
		return nil, err
	}

	return spec.UpdateFollowRequest200Response{}, nil
}
