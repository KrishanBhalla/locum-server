package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/KrishanBhalla/locum-server/api/spec"
	"github.com/KrishanBhalla/locum-server/models"
)

func UpdateLocation(ctx context.Context, request spec.UpdateLocationRequestObject) (spec.UpdateLocationResponseObject, error) {

	internalServerError := spec.UpdateLocationdefaultResponse{StatusCode: http.StatusInternalServerError}

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

	userLocationService := services.UserLocation

	err = userLocationService.Append(models.UserLocation{
		UserId: userId,
		GeoTimes: []models.GeoTime{
			{
				Latitude:  request.Body.Latitude,
				Longitude: request.Body.Longitude,
				Timestamp: time.UnixMilli(request.Body.Timestamp),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return internalServerError, err
	}

	return spec.UpdateLocation200Response{}, nil
}
