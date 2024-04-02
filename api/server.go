package api

import (
	"context"

	"github.com/KrishanBhalla/locum-server/api/routes"
	"github.com/KrishanBhalla/locum-server/api/spec"
)

type ServerImpl struct {
}

// FindFriendRequests implements spec.StrictServerInterface.
func (*ServerImpl) GetFriendRequests(ctx context.Context, request spec.GetFriendRequestsRequestObject) (spec.GetFriendRequestsResponseObject, error) {
	return routes.GetFriendRequests(ctx, request)
}

// CreateFriendRequest implements spec.StrictServerInterface.
func (*ServerImpl) CreateFriendRequest(ctx context.Context, request spec.CreateFriendRequestRequestObject) (spec.CreateFriendRequestResponseObject, error) {
	return routes.CreateFriendRequest(ctx, request)
}

// UpdateFriendRequest implements spec.StrictServerInterface.
func (*ServerImpl) UpdateFriendRequest(ctx context.Context, request spec.UpdateFriendRequestRequestObject) (spec.UpdateFriendRequestResponseObject, error) {
	return routes.UpdateFriendRequest(ctx, request)
}

// UpdateLocationsOfFriendedUsers implements spec.StrictServerInterface.
func (*ServerImpl) GetLocationsOfFriends(ctx context.Context, request spec.GetLocationsOfFriendsRequestObject) (spec.GetLocationsOfFriendsResponseObject, error) {
	return routes.GetLocationsOfFriends(ctx)
}

// DeleteFriend implements spec.StrictServerInterface.
func (*ServerImpl) DeleteFriend(ctx context.Context, request spec.DeleteFriendRequestObject) (spec.DeleteFriendResponseObject, error) {
	return routes.DeleteFriend(ctx, request)
}

// FindFriends implements spec.StrictServerInterface.
func (*ServerImpl) GetFriends(ctx context.Context, request spec.GetFriendsRequestObject) (spec.GetFriendsResponseObject, error) {
	return routes.GetFriends(ctx, request)
}

// FindUsers implements spec.StrictServerInterface.
func (*ServerImpl) FindUsers(ctx context.Context, request spec.FindUsersRequestObject) (spec.FindUsersResponseObject, error) {
	return routes.FindUsers(ctx, request)
}

// LoginOrSignup implements spec.StrictServerInterface.
func (*ServerImpl) LoginOrSignup(ctx context.Context, request spec.LoginOrSignupRequestObject) (spec.LoginOrSignupResponseObject, error) {
	return routes.SignupOrLogin(ctx, request)
}

var _ spec.StrictServerInterface = &ServerImpl{}
