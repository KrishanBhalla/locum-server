package routes

import (
	"encoding/json"
	"net/http"

	"github.com/KrishanBhalla/locum-server/services"
)

type FindFriendsRequest struct {
	UserId string `json:"userId"`
}
type Friend struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`
}
type FindFriendsResponse struct {
	Followers []Friend `json:"followers"`
	Following []Friend `json:"following"`
}

func FindFriends(services services.Services, w http.ResponseWriter, r *http.Request) {
	// Decode body
	var req FindFriendsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Find friends
	friends, err := services.UserFriends.ByUserID(req.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Create Response
	followers := make([]Friend, 0, len(friends.FollowerUserIds))
	for _, follower := range friends.FollowerUserIds {
		user, err := services.User.ByID(follower)
		if err == nil {
			followers = append(followers, Friend{UserId: follower, Name: user.FullName})
		}
	}
	followed := make([]Friend, 0, len(friends.FollowingUserIds))
	for _, following := range friends.FollowingUserIds {
		user, err := services.User.ByID(following)
		if err == nil {
			followed = append(followed, Friend{UserId: following, Name: user.FullName})
		}
	}
	// Write Response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	err = json.NewEncoder(w).Encode(&FindFriendsResponse{Followers: followers, Following: followed})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
