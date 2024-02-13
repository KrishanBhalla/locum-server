package api

import (
	"context"

	"github.com/KrishanBhalla/locum-server/api/routes"
	"github.com/KrishanBhalla/locum-server/api/spec"
)

type ServerImpl struct {
}

// DeleteFollower implements spec.StrictServerInterface.
func (*ServerImpl) DeleteFollower(ctx context.Context, request spec.DeleteFollowerRequestObject) (spec.DeleteFollowerResponseObject, error) {
	panic("unimplemented")
}

// DeleteFollowing implements spec.StrictServerInterface.
func (*ServerImpl) DeleteFollowing(ctx context.Context, request spec.DeleteFollowingRequestObject) (spec.DeleteFollowingResponseObject, error) {
	panic("unimplemented")
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
	panic("unimplemented")
}

// LoginOrSignup implements spec.StrictServerInterface.
func (*ServerImpl) LoginOrSignup(ctx context.Context, request spec.LoginOrSignupRequestObject) (spec.LoginOrSignupResponseObject, error) {
	return routes.SignupOrLogin(ctx, request)
}

var _ spec.StrictServerInterface = &ServerImpl{}
