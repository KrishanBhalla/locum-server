package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/KrishanBhalla/locum-server/api/spec"
	"github.com/KrishanBhalla/locum-server/models"
)

func CreateFriendRequest(ctx context.Context, request spec.CreateFriendRequestRequestObject) (spec.CreateFriendRequestResponseObject, error) {

	internalServerError := spec.CreateFriendRequestdefaultResponse{StatusCode: http.StatusInternalServerError}

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
	userToBefriend := request.Body.FriendId
	reqTime := time.Now()
	userFriends := services.UserFriends
	err = userFriends.AddFriendRequest(userToBefriend, models.FriendRequest{UserId: userId, Timestamp: reqTime})
	if err != nil {
		return nil, err
	}
	err = userFriends.AddFriendRequest(userId, models.FriendRequest{UserId: userToBefriend, Timestamp: reqTime})
	if err != nil {
		return nil, err
	}

	return spec.CreateFriendRequest200Response{}, nil
}
