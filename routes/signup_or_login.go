package routes

import (
	"encoding/json"
	"net/http"

	"github.com/KrishanBhalla/locum-server/models"
	"github.com/KrishanBhalla/locum-server/services"
	"github.com/dgraph-io/badger"
)

type LoginRequest struct {
	UserId   string `json:"userId"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
}

func SignupOrLogin(services services.Services, w http.ResponseWriter, r *http.Request) {
	// Decode Request
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userService := services.User
	user, err := userService.ByID(req.UserId)
	if err != nil && err != badger.ErrKeyNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if err == badger.ErrKeyNotFound {
		err = userService.Create(models.User{Id: req.UserId, FullName: req.FullName, Email: req.Email})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := userService.Update(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.WriteHeader(http.StatusOK)
}
