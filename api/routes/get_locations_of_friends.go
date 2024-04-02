package routes

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/KrishanBhalla/locum-server/api/spec"
	chiMw "github.com/go-chi/chi/middleware"
)

func GetLocationsOfFriends(ctx context.Context) (spec.GetLocationsOfFriendsResponseObject, error) {
	/// setup
	reqId := chiMw.GetReqID(ctx)
	internalServerError := spec.GetLocationsOfFriendsdefaultResponse{StatusCode: http.StatusInternalServerError}

	services, err := validateServices(ctx)
	if err != nil {
		return internalServerError, err
	}

	userToken, err := validateToken(ctx, services)
	if err != nil {
		return spec.UnauthorizedErrorResponse{}, err
	}

	// Process
	userFriendsService := services.UserFriends
	userLocationService := services.UserLocation

	userFriends, err := userFriendsService.ByUserID(userToken.UserId)
	if err != nil {
		return internalServerError, errors.New(fmt.Sprintf("Failed to find user friends with reqId: %s, err: %s", reqId, err.Error()))
	}
	friends := userFriends.FriendIds

	locations := make(spec.GetLocationsOfFriends200JSONResponse, 0)
	for _, f := range friends {
		loc, err := userLocationService.LatestGeoTimeByUserID(f)
		if err != nil {
			log.Default().Print(err)
		}
		locations = append(locations, spec.UserLocation{UserId: f, Latitude: loc.Latitude, Longitude: loc.Longitude, Timestamp: loc.Timestamp.Unix()})
	}
	return locations, nil
}
