package routes

import (
	"context"
	"net/http"

	"github.com/KrishanBhalla/locum-server/api/spec"
)

func RenameUser(ctx context.Context, request spec.RenameUserRequestObject) (spec.RenameUserResponseObject, error) {

	internalServerError := spec.RenameUserdefaultResponse{StatusCode: http.StatusInternalServerError}

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
	newName := request.Body.NewName

	userService := services.User

	user, err := userService.ByID(userId)
	if err != nil {
		return internalServerError, err
	}
	user.FullName = newName
	err = userService.Update(user)
	if err != nil {
		return internalServerError, err
	}

	return spec.RenameUser200Response{}, nil
}
