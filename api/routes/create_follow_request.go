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

func CreateFollowRequest(ctx context.Context, request spec.CreateFollowRequestRequestObject) (spec.CreateFollowRequestResponseObject, error) {

	services, ok := services.FromContext(ctx)
	reqId := chiMw.GetReqID(ctx)
	internalServerError := spec.CreateFollowRequestdefaultResponse{StatusCode: http.StatusInternalServerError}
	if !ok {
		return internalServerError, errors.New(fmt.Sprintf("No services passed via context, reqId: %s", reqId))
	}

	// Process
	userId := request.Body.UserId
	userToFollow := request.Body.FollowingUserId
	reqTime := time.Now()
	userFriends := services.UserFriends
	err := userFriends.AddFollowerRequest(userToFollow, models.FollowRequest{UserId: userId, Timestamp: reqTime})
	if err != nil {
		return nil, err
	}
	err = userFriends.AddFollowRequest(userId, models.FollowRequest{UserId: userToFollow, Timestamp: reqTime})
	if err != nil {
		return nil, err
	}

	return spec.CreateFollowRequest200Response{}, nil
}
