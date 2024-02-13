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

func FindFollowers(ctx context.Context, request spec.FindFollowersRequestObject) (spec.FindFollowersResponseObject, error) {
	/// setup
	services, ok := services.FromContext(ctx)
	reqId := chiMw.GetReqID(ctx)
	internalServerError := spec.FindFollowersdefaultResponse{StatusCode: http.StatusInternalServerError}
	if !ok {
		return internalServerError, errors.New(fmt.Sprintf("No services passed via context, reqId: %s", reqId))
	}

	// Process

	friends, err := services.UserFriends.ByUserID(request.Body.UserId)
	if err != nil {
		return internalServerError, err
	}
	// Create Response
	followers := make(spec.FindFollowers200JSONResponse, 0, len(friends.FollowerUserIds))
	for _, follower := range friends.FollowerUserIds {
		user, err := services.User.ByID(follower)
		if err == nil {
			followers = append(followers, spec.UserResponse{UserId: follower, FullName: user.FullName})
		}
	}
	return followers, nil
}
