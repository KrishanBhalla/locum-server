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

func UpdateFriendRequest(ctx context.Context, request spec.UpdateFriendRequestRequestObject) (spec.UpdateFriendRequestResponseObject, error) {

	services, ok := services.FromContext(ctx)
	reqId := chiMw.GetReqID(ctx)
	internalServerError := spec.UpdateFriendRequestdefaultResponse{StatusCode: http.StatusInternalServerError}
	if !ok {
		return internalServerError, errors.New(fmt.Sprintf("No services passed via context, reqId: %s", reqId))
	}

	// Process
	userId := request.Body.UserId
	requestingUser := request.Body.FriendId
	requestAccepted := request.Body.Accept

	userFriends := services.UserFriends
	var err error
	if requestAccepted {
		err = userFriends.AddFriend(userId, requestingUser)
	}
	if err != nil {
		return nil, err
	}
	err = userFriends.RemoveFriendRequest(userId, requestingUser)
	if err != nil {
		return nil, err
	}

	return spec.UpdateFriendRequest200Response{}, nil
}
