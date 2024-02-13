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

func DeleteFollower(ctx context.Context, request spec.DeleteFollowerRequestObject) (spec.DeleteFollowerResponseObject, error) {
	/// setup
	services, ok := services.FromContext(ctx)
	reqId := chiMw.GetReqID(ctx)
	internalServerError := spec.DeleteFollowerdefaultResponse{StatusCode: http.StatusInternalServerError}
	if !ok {
		return internalServerError, errors.New(fmt.Sprintf("No services passed via context, reqId: %s", reqId))
	}

	// Process
	// the userId is followed by followerUserId, so we must Remove the followerUserId from the userId's followers
	// Conversely the followerUserId is following the userId, so we must Remove the userId from the followerUserId's following group
	userId := request.Body.UserId
	followerUserId := request.Body.FollowerUserId

	err := services.UserFriends.RemoveFollower(userId, followerUserId)
	if err != nil {
		return internalServerError, err
	}
	err = services.UserFriends.RemoveFollowing(followerUserId, userId)
	if err != nil {
		return internalServerError, err
	}
	return spec.DeleteFollower204Response{}, nil
}
