package routes

import (
	"context"
	"net/http"

	"github.com/KrishanBhalla/locum-server/api/spec"
)

func DeleteFriend(ctx context.Context, request spec.DeleteFriendRequestObject) (spec.DeleteFriendResponseObject, error) {
	/// setup

	internalServerError := spec.DeleteFrienddefaultResponse{StatusCode: http.StatusInternalServerError}

	services, err := validateServices(ctx)
	if err != nil {
		return internalServerError, err
	}

	userToken, err := validateToken(ctx, services)
	if err != nil {
		return spec.UnauthorizedErrorResponse{}, err
	}

	// Process
	// the userId is followed by followerUserId, so we must Remove the followerUserId from the userId's followers
	// Conversely the followerUserId is following the userId, so we must Remove the userId from the followerUserId's following group
	userId := userToken.UserId
	friendId := request.Body.FriendId

	err = services.UserFriends.RemoveFriend(userId, friendId)
	if err != nil {
		return internalServerError, err
	}
	return spec.DeleteFriend204Response{}, nil
}
