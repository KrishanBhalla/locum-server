package models

import (
	"encoding/json"
	"time"

	"github.com/dgraph-io/badger"
)

type FriendRequest struct {
	UserId    string    `json:"userId"`
	Timestamp time.Time `json:"timestamp"`
}

func concatenateFriendRequests(initialRequests []FriendRequest, newRequests ...FriendRequest) []FriendRequest {

	requestMap := make(map[string]time.Time)
	for _, f := range initialRequests {
		requestMap[f.UserId] = f.Timestamp
	}

	for _, followReq := range newRequests {
		reqTs, ok := requestMap[followReq.UserId]
		if !ok || (ok && followReq.Timestamp.Before(reqTs)) {
			requestMap[followReq.UserId] = followReq.Timestamp
		}
	}
	result := make([]FriendRequest, 0, len(requestMap))
	for k, v := range requestMap {
		result = append(result, FriendRequest{UserId: k, Timestamp: v})
	}
	return result
}

func removeFriendRequest(original []FriendRequest, key string) []FriendRequest {
	for i, v := range original {
		if v.UserId == key {
			if i < len(original) {
				return append(original[:i], original[i+1:]...)
			} else {
				return original[:len(original)-1]
			}
		}
	}
	return original
}

type UserFriends struct {
	UserId         string          `json:"id"`
	FriendIds      []string        `json:"friendIds"`
	FriendRequests []FriendRequest `json:"followRequests"`
}

type UserFriendsDB interface {
	ByUserID(userId string) (UserFriends, error)

	// Methods for altering contents
	Create(userFriends UserFriends) error
	Append(userFriends UserFriends) error

	RemoveFriend(userId, friendId string) error
	AddFriend(userId, friendId string) error

	RemoveFriendRequest(userId, friendId string) error
	AddFriendRequest(userId string, followRequest FriendRequest) error

	Update(userFriends UserFriends) error
	Delete(userId string) error
	DbCloser
}

// Define userFriendsDB and ensure it implements UserFriendsDB
var _ UserFriendsDB = &userFriendsDB{}

type userFriendsDB struct {
	db *badger.DB
}

func NewUserFriendsDB(db *badger.DB) UserFriendsDB {
	return &userFriendsDB{db}
}

// ByID implements UserFriendsDB.
func (db *userFriendsDB) ByUserID(userId string) (UserFriends, error) {

	var userFriends UserFriends
	var data = make([]byte, 0)
	data, err := lookupByKey(db.db, userId, data)
	if err != nil {
		return UserFriends{}, err
	}
	err = json.Unmarshal(data, &userFriends)
	if err != nil {
		return UserFriends{}, err
	}
	return userFriends, nil
}

// Create implements UserFriendsDB.
func (db *userFriendsDB) Create(userFriends UserFriends) error {
	return db.Update(userFriends)
}

// Delete implements UserFriendsDB.
func (db *userFriendsDB) Delete(userId string) error {
	userFriends, err := db.ByUserID(userId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	if err == badger.ErrKeyNotFound {
		return nil
	}

	err = db.db.Update(func(txn *badger.Txn) error {

		userFriendsBytes, err := json.Marshal(userFriends)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(userFriends.UserId), userFriendsBytes)
		return err
	})
	return err
}

func (db *userFriendsDB) Append(userFriends UserFriends) error {

	friends, err := db.ByUserID(userFriends.UserId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	userFriends.FriendIds = appendToSliceWithoutDuplicates(friends.FriendIds, userFriends.FriendIds...)
	userFriends.FriendRequests = concatenateFriendRequests(userFriends.FriendRequests, friends.FriendRequests...)

	return db.Update(userFriends)
}

func (db *userFriendsDB) RemoveFriend(userId, friendId string) error {

	userFriends, err := db.ByUserID(userId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	userFriends.FriendIds = removeFromSlice(userFriends.FriendIds, friendId)
	return db.Update(userFriends)
}

func (db *userFriendsDB) AddFriend(userId, friendId string) error {

	userFriends, err := db.ByUserID(userId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	userFriends.FriendIds = appendToSliceWithoutDuplicates(userFriends.FriendIds, friendId)
	return db.Update(userFriends)
}

// Follow reqeusts

func (db *userFriendsDB) RemoveFriendRequest(userId, friendRequestUserId string) error {

	userFriends, err := db.ByUserID(userId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	userFriends.FriendRequests = removeFriendRequest(userFriends.FriendRequests, friendRequestUserId)
	return db.Update(userFriends)
}

func (db *userFriendsDB) AddFriendRequest(userId string, friendRequest FriendRequest) error {

	userFriends, err := db.ByUserID(userId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	userFriends.FriendRequests = concatenateFriendRequests(userFriends.FriendRequests, friendRequest)
	return db.Update(userFriends)
}

// Update implements UserFriendsDB.
func (db *userFriendsDB) Update(userFriends UserFriends) error {

	err := db.db.Update(func(txn *badger.Txn) error {

		userFriendsBytes, err := json.Marshal(userFriends)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(userFriends.UserId), userFriendsBytes)
		return err
	})
	return err
}

func (db *userFriendsDB) CloseDB() error {
	return db.db.Close()
}
