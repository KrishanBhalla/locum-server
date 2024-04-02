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

func GetFriends(ctx context.Context, request spec.GetFriendsRequestObject) (spec.GetFriendsResponseObject, error) {
	/// setup
	services, ok := services.FromContext(ctx)
	reqId := chiMw.GetReqID(ctx)
	internalServerError := spec.GetFriendsdefaultResponse{StatusCode: http.StatusInternalServerError}
	if !ok {
		return internalServerError, errors.New(fmt.Sprintf("No services passed via context, reqId: %s", reqId))
	}

	// Process

	friends, err := services.UserFriends.ByUserID(request.Body.UserId)
	if err != nil {
		return internalServerError, err
	}
	// Create Response
	followers := make(spec.GetFriends200JSONResponse, 0, len(friends.FriendIds))
	for _, follower := range friends.FriendIds {
		user, err := services.User.ByID(follower)
		if err == nil {
			followers = append(followers, spec.UserResponse{UserId: follower, FullName: user.FullName})
		}
	}
	return followers, nil
}
