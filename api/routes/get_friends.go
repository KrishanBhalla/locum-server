package routes

import (
	"context"
	"net/http"

	"github.com/KrishanBhalla/locum-server/api/spec"
	"github.com/dgraph-io/badger"
)

func GetFriends(ctx context.Context, request spec.GetFriendsRequestObject) (spec.GetFriendsResponseObject, error) {
	/// setup
	internalServerError := spec.GetFriendsdefaultResponse{StatusCode: http.StatusInternalServerError}

	services, err := validateServices(ctx)
	if err != nil {
		return internalServerError, err
	}

	userToken, err := validateToken(ctx, services)
	if err != nil {
		return spec.UnauthorizedErrorResponse{}, err
	}

	// Process

	friends, err := services.UserFriends.ByUserID(userToken.UserId)
	if err != nil && err != badger.ErrKeyNotFound {
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
