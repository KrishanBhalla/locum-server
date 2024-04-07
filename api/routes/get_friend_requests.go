package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/KrishanBhalla/locum-server/api/spec"
)

func GetFriendRequests(ctx context.Context, request spec.GetFriendRequestsRequestObject) (spec.GetFriendRequestsResponseObject, error) {

	internalServerError := spec.GetFriendRequestsdefaultResponse{StatusCode: http.StatusInternalServerError}

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
	userFriendsService := services.UserFriends
	userService := services.User
	friends, err := userFriendsService.ByUserID(userId)
	if err != nil {
		return nil, err
	}

	friendRequests := friends.FriendRequests
	response := make(spec.GetFriendRequests200JSONResponse, 0, len(friendRequests))
	for _, f := range friendRequests {

		follower, err := userService.ByID(f.UserId)
		if err != nil {
			log.Default().Printf("Failed to find user for id %s in find_follower_requests.go. Err %s", f.UserId, err.Error())
		} else {
			response = append(response, spec.UserResponse{UserId: follower.Id, FullName: follower.FullName})
		}
	}

	return response, nil
}
