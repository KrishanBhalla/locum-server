package routes

import (
	"context"
	"errors"
	"fmt"

	"github.com/KrishanBhalla/locum-server/models"
	"github.com/KrishanBhalla/locum-server/services"
	"github.com/KrishanBhalla/locum-server/services/tokens"
	chiMw "github.com/go-chi/chi/v5/middleware"
)

func validateServices(ctx context.Context) (*services.Services, error) {
	services, ok := services.FromContext(ctx)
	reqId := chiMw.GetReqID(ctx)

	if !ok {
		return nil, errors.New(fmt.Sprintf("No services passed via context, reqId: %s", reqId))
	}
	return services, nil
}

func validateToken(ctx context.Context, services *services.Services) (*models.UserToken, error) {
	token, authError := tokens.FromContext(ctx)
	reqId := chiMw.GetReqID(ctx)
	fmt.Println("CONTEXT")

	if authError != nil {
		fmt.Println("Auth Error", authError.Error())
		return nil, authError
	}

	userToken, err := services.UserToken.ByToken(token.Token)
	fmt.Println("HAVE USER TOKEN")
	if err != nil {
		fmt.Println("Err", err.Error())
		return nil, errors.New(fmt.Sprintf("reqId: %s. Error %s", reqId, err.Error()))
	}

	return &userToken, nil
}
