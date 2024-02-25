package models

import (
	"encoding/json"
	"time"

	"github.com/dgraph-io/badger"
)

type FollowRequest struct {
	UserId    string    `json:"userId"`
	Timestamp time.Time `json:"timestamp"`
}

func concatenateFollowRequests(initialRequests []FollowRequest, newRequests ...FollowRequest) []FollowRequest {

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
	result := make([]FollowRequest, 0, len(requestMap))
	for k, v := range requestMap {
		result = append(result, FollowRequest{UserId: k, Timestamp: v})
	}
	return result
}

func removeFollowRequest(original []FollowRequest, key string) []FollowRequest {
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
	UserId           string          `json:"id"`
	FollowerUserIds  []string        `json:"followerUserIds"`
	FollowingUserIds []string        `json:"followingUserIds"`
	FollowerRequests []FollowRequest `json:"followerRequests"`
	FollowRequests   []FollowRequest `json:"followRequests"`
}

type UserFriendsDB interface {
	ByUserID(userId string) (UserFriends, error)

	// Methods for altering contents
	Create(userFriends UserFriends) error
	Append(userFriends UserFriends) error

	RemoveFollower(userId, followerUserId string) error
	RemoveFollowing(userId, followingUserId string) error
	AddFollower(userId, followerUserId string) error
	AddFollowing(userId, followingUserId string) error

	RemoveFollowerRequest(userId, followerRequestUserId string) error
	RemoveFollowRequest(userId, followRequestUserId string) error
	AddFollowerRequest(userId string, followerRequest FollowRequest) error
	AddFollowRequest(userId string, followRequest FollowRequest) error

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

	userFriends.FollowerUserIds = appendToSliceWithoutDuplicates(friends.FollowerUserIds, userFriends.FollowerUserIds...)
	userFriends.FollowingUserIds = appendToSliceWithoutDuplicates(friends.FollowingUserIds, userFriends.FollowingUserIds...)

	userFriends.FollowRequests = concatenateFollowRequests(userFriends.FollowRequests, friends.FollowRequests...)
	userFriends.FollowerRequests = concatenateFollowRequests(userFriends.FollowerRequests, friends.FollowerRequests...)

	return db.Update(userFriends)
}

func (db *userFriendsDB) RemoveFollower(userId, followerUserId string) error {

	userFriends, err := db.ByUserID(userId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	userFriends.FollowerUserIds = removeFromSlice(userFriends.FollowerUserIds, followerUserId)
	return db.Update(userFriends)
}
func (db *userFriendsDB) RemoveFollowing(userId, followingUserId string) error {

	userFriends, err := db.ByUserID(userId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	userFriends.FollowingUserIds = removeFromSlice(userFriends.FollowingUserIds, followingUserId)

	return db.Update(userFriends)
}

func (db *userFriendsDB) AddFollower(userId, followerUserId string) error {

	userFriends, err := db.ByUserID(userId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	userFriends.FollowerUserIds = appendToSliceWithoutDuplicates(userFriends.FollowerUserIds, followerUserId)
	return db.Update(userFriends)
}

func (db *userFriendsDB) AddFollowing(userId, followingUserId string) error {

	userFriends, err := db.ByUserID(userId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	userFriends.FollowingUserIds = appendToSliceWithoutDuplicates(userFriends.FollowingUserIds, followingUserId)

	return db.Update(userFriends)
}

// Follow reqeusts

func (db *userFriendsDB) RemoveFollowerRequest(userId, followerRequestUserId string) error {

	userFriends, err := db.ByUserID(userId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	userFriends.FollowRequests = removeFollowRequest(userFriends.FollowRequests, followerRequestUserId)
	return db.Update(userFriends)
}
func (db *userFriendsDB) RemoveFollowRequest(userId, followRequestUserId string) error {

	userFriends, err := db.ByUserID(userId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	userFriends.FollowRequests = removeFollowRequest(userFriends.FollowRequests, followRequestUserId)

	return db.Update(userFriends)
}

func (db *userFriendsDB) AddFollowerRequest(userId string, followerRequest FollowRequest) error {

	userFriends, err := db.ByUserID(userId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	userFriends.FollowerRequests = concatenateFollowRequests(userFriends.FollowerRequests, followerRequest)
	return db.Update(userFriends)
}

func (db *userFriendsDB) AddFollowRequest(userId string, followRequest FollowRequest) error {

	userFriends, err := db.ByUserID(userId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	userFriends.FollowRequests = concatenateFollowRequests(userFriends.FollowRequests, followRequest)

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
