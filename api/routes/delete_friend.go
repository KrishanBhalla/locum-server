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

func DeleteFriend(ctx context.Context, request spec.DeleteFriendRequestObject) (spec.DeleteFriendResponseObject, error) {
	/// setup
	services, ok := services.FromContext(ctx)
	reqId := chiMw.GetReqID(ctx)
	internalServerError := spec.DeleteFrienddefaultResponse{StatusCode: http.StatusInternalServerError}
	if !ok {
		return internalServerError, errors.New(fmt.Sprintf("No services passed via context, reqId: %s", reqId))
	}

	// Process
	// the userId is followed by followerUserId, so we must Remove the followerUserId from the userId's followers
	// Conversely the followerUserId is following the userId, so we must Remove the userId from the followerUserId's following group
	userId := request.Body.UserId
	friendId := request.Body.FriendId

	err := services.UserFriends.RemoveFriend(userId, friendId)
	if err != nil {
		return internalServerError, err
	}
	return spec.DeleteFriend204Response{}, nil
}
