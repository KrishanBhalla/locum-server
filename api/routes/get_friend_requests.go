package routes

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/KrishanBhalla/locum-server/api/spec"
	"github.com/KrishanBhalla/locum-server/services"
	chiMw "github.com/go-chi/chi/middleware"
)

func GetFriendRequests(ctx context.Context, request spec.GetFriendRequestsRequestObject) (spec.GetFriendRequestsResponseObject, error) {

	services, ok := services.FromContext(ctx)
	reqId := chiMw.GetReqID(ctx)
	internalServerError := spec.GetFriendRequestsdefaultResponse{StatusCode: http.StatusInternalServerError}
	if !ok {
		return internalServerError, errors.New(fmt.Sprintf("No services passed via context, reqId: %s", reqId))
	}

	// Process
	userId := request.Body.UserId
	userFriendsService := services.UserFriends
	userService := services.User
	friends, err := userFriendsService.ByUserID(userId)
	if err != nil {
		return nil, err
	}

	followerRequests := friends.FriendRequests
	response := make([]spec.UserResponse, 0, len(followerRequests))
	for _, f := range followerRequests {

		follower, err := userService.ByID(f.UserId)
		if err != nil {
			log.Default().Printf("Failed to find user for id %s in find_follower_requests.go. Err %s", f.UserId, err.Error())
		} else {
			response = append(response, spec.UserResponse{UserId: follower.Id, FullName: follower.FullName})
		}
	}

	type UserResponse struct {
		FullName string `json:"fullName"`
		UserId   string `json:"userId"`
	}

	return spec.GetFriendRequests200JSONResponse{}, nil
}
