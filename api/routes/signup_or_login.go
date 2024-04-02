package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/KrishanBhalla/locum-server/api/spec"
	"github.com/KrishanBhalla/locum-server/models"
	"github.com/dgraph-io/badger"
	chiMw "github.com/go-chi/chi/middleware"
)

func SignupOrLogin(ctx context.Context, request spec.LoginOrSignupRequestObject) (spec.LoginOrSignupResponseObject, error) {
	// setup
	req := request.Body

	services, err := validateServices(ctx)
	reqId := chiMw.GetReqID(ctx)
	internalServerError := spec.LoginOrSignupdefaultResponse{StatusCode: http.StatusInternalServerError}
	if err != nil {
		return internalServerError, err
	}
	userService := services.User

	// Process
	user, err := userService.ByID(req.UserId)
	if err != nil && err != badger.ErrKeyNotFound {
		return nil, err
	} else if err == badger.ErrKeyNotFound {
		defaultNilValue := "Unknown"
		if req.FullName == nil {
			req.FullName = &defaultNilValue
		}
		if req.Email == nil {
			req.FullName = &defaultNilValue
		}
		err := userService.Create(models.User{Id: req.UserId, FullName: *req.FullName, Email: *req.Email})
		if err != nil {
			return internalServerError, errors.New(fmt.Sprintf("Error creating user (SignupOrLogin) %s", reqId))
		}

	} else {
		err := userService.Update(user)
		if err != nil {
			return internalServerError, errors.New(fmt.Sprintf("Error updating user (SignupOrLogin) %s", reqId))
		}
		// now updates are done,m get the user again
		user, err = userService.ByID(user.Id)
		if err != nil {
			return internalServerError, err
		}
	}

	tokenService := services.UserToken
	token := models.NewUserToken(req.UserId)
	err = tokenService.Create(token)
	if err != nil {
		return internalServerError, err
	}

	return spec.LoginOrSignup200JSONResponse{Token: "BEARER " + token.Token}, nil
}
