package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/KrishanBhalla/locum-server/api/spec"
	"github.com/KrishanBhalla/locum-server/models"
	"github.com/KrishanBhalla/locum-server/services"
	"github.com/dgraph-io/badger"
	chiMw "github.com/go-chi/chi/middleware"
)

func SignupOrLogin(ctx context.Context, request spec.LoginOrSignupRequestObject) (spec.LoginOrSignupResponseObject, error) {
	// setup
	req := request.Body

	services, ok := services.FromContext(ctx)
	reqId := chiMw.GetReqID(ctx)
	internalServerError := spec.LoginOrSignupdefaultResponse{StatusCode: http.StatusInternalServerError}
	if !ok {
		return internalServerError, errors.New(fmt.Sprintf("No services passed via context, reqId: %s", reqId))
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
		err = userService.Create(models.User{Id: req.UserId, FullName: *req.FullName, Email: *req.Email})
		if err != nil {
			return internalServerError, errors.New(fmt.Sprintf("Error creating user (SignupOrLogin) %s", reqId))
		}
	} else {
		err := userService.Update(user)
		if err != nil {
			return internalServerError, errors.New(fmt.Sprintf("Error updating user (SignupOrLogin) %s", reqId))
		}
	}
	return spec.LoginOrSignup200Response{}, nil
}
