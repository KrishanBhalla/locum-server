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

func DeleteFollowing(ctx context.Context, request spec.DeleteFollowingRequestObject) (spec.DeleteFollowingResponseObject, error) {
	/// setup
	services, ok := services.FromContext(ctx)
	reqId := chiMw.GetReqID(ctx)
	internalServerError := spec.DeleteFollowingdefaultResponse{StatusCode: http.StatusInternalServerError}
	if !ok {
		return internalServerError, errors.New(fmt.Sprintf("No services passed via context, reqId: %s", reqId))
	}

	// Process
	// the userId is following followingUserId, so we must Remove the followerUserId from the userId's following group
	// Conversely the followingUserId is followed by the userId, so we must Remove the userId from the followingUserId's followers group
	userId := request.Body.UserId
	followingUserId := request.Body.FollowingUserId

	err := services.UserFriends.RemoveFollowing(userId, followingUserId)
	if err != nil {
		return internalServerError, err
	}
	err = services.UserFriends.RemoveFollower(followingUserId, userId)
	if err != nil {
		return internalServerError, err
	}
	return spec.DeleteFollowing204Response{}, nil
}
