package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/KrishanBhalla/locum-server/api/spec"
	"github.com/KrishanBhalla/locum-server/models"
	"github.com/KrishanBhalla/locum-server/services"
	chiMw "github.com/go-chi/chi/middleware"
)

func CreateFriendRequest(ctx context.Context, request spec.CreateFriendRequestRequestObject) (spec.CreateFriendRequestResponseObject, error) {

	services, ok := services.FromContext(ctx)
	reqId := chiMw.GetReqID(ctx)
	internalServerError := spec.CreateFriendRequestdefaultResponse{StatusCode: http.StatusInternalServerError}
	if !ok {
		return internalServerError, errors.New(fmt.Sprintf("No services passed via context, reqId: %s", reqId))
	}

	// Process
	userId := request.Body.UserId
	userToBefriend := request.Body.FriendId
	reqTime := time.Now()
	userFriends := services.UserFriends
	err := userFriends.AddFriendRequest(userToBefriend, models.FriendRequest{UserId: userId, Timestamp: reqTime})
	if err != nil {
		return nil, err
	}
	err = userFriends.AddFriendRequest(userId, models.FriendRequest{UserId: userToBefriend, Timestamp: reqTime})
	if err != nil {
		return nil, err
	}

	return spec.CreateFriendRequest200Response{}, nil
}
