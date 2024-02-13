package api

import (
	"context"

	"github.com/KrishanBhalla/locum-server/api/routes"
	"github.com/KrishanBhalla/locum-server/api/spec"
)

type ServerImpl struct {
}

// CreateFollowRequest implements spec.StrictServerInterface.
func (*ServerImpl) CreateFollowRequest(ctx context.Context, request spec.CreateFollowRequestRequestObject) (spec.CreateFollowRequestResponseObject, error) {
	panic("unimplemented")
}

// UpdateFollowRequest implements spec.StrictServerInterface.
func (*ServerImpl) UpdateFollowRequest(ctx context.Context, request spec.UpdateFollowRequestRequestObject) (spec.UpdateFollowRequestResponseObject, error) {
	panic("unimplemented")
}

// UpdateLocationsOfFollowedUsers implements spec.StrictServerInterface.
func (*ServerImpl) UpdateLocationsOfFollowedUsers(ctx context.Context, request spec.UpdateLocationsOfFollowedUsersRequestObject) (spec.UpdateLocationsOfFollowedUsersResponseObject, error) {
	panic("unimplemented")
}

// DeleteFollower implements spec.StrictServerInterface.
func (*ServerImpl) DeleteFollower(ctx context.Context, request spec.DeleteFollowerRequestObject) (spec.DeleteFollowerResponseObject, error) {
	return routes.DeleteFollower(ctx, request)
}

// DeleteFollowing implements spec.StrictServerInterface.
func (*ServerImpl) DeleteFollowing(ctx context.Context, request spec.DeleteFollowingRequestObject) (spec.DeleteFollowingResponseObject, error) {
	return routes.DeleteFollowing(ctx, request)
}

// FindFollowers implements spec.StrictServerInterface.
func (*ServerImpl) FindFollowers(ctx context.Context, request spec.FindFollowersRequestObject) (spec.FindFollowersResponseObject, error) {
	return routes.FindFollowers(ctx, request)
}

// FindFollowing implements spec.StrictServerInterface.
func (*ServerImpl) FindFollowing(ctx context.Context, request spec.FindFollowingRequestObject) (spec.FindFollowingResponseObject, error) {
	return routes.FindFollowing(ctx, request)
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
