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

func FindFollowing(ctx context.Context, request spec.FindFollowingRequestObject) (spec.FindFollowingResponseObject, error) {
	/// setup
	req := request.Body

	services, ok := services.FromContext(ctx)
	reqId := chiMw.GetReqID(ctx)
	internalServerError := spec.FindFollowingdefaultResponse{StatusCode: http.StatusInternalServerError}
	if !ok {
		return internalServerError, errors.New(fmt.Sprintf("No services passed via context, reqId: %s", reqId))
	}

	// Process

	friends, err := services.UserFriends.ByUserID(req.UserId)
	if err != nil {
		return internalServerError, err
	}
	// Create Response
	followingUsers := make(spec.FindFollowing200JSONResponse, 0, len(friends.FollowerUserIds))
	for _, following := range friends.FollowingUserIds {
		user, err := services.User.ByID(following)
		if err == nil {
			followingUsers = append(followingUsers, spec.UserResponse{UserId: following, FullName: user.FullName})
		}
	}
	return followingUsers, nil
}
