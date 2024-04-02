package routes

import (
	"context"
	"net/http"

	"github.com/KrishanBhalla/locum-server/api/spec"
)

func UpdateFriendRequest(ctx context.Context, request spec.UpdateFriendRequestRequestObject) (spec.UpdateFriendRequestResponseObject, error) {

	internalServerError := spec.UpdateFriendRequestdefaultResponse{StatusCode: http.StatusInternalServerError}

	services, err := validateServices(ctx)
	if err != nil {
		return internalServerError, err
	}

	userToken, err := validateToken(ctx, services)
	if err != nil {
		return spec.UnauthorizedErrorResponse{}, err
	}

	// Process
	userId := userToken.UserId
	requestingUser := request.Body.FriendId
	requestAccepted := request.Body.Accept

	userFriends := services.UserFriends

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
