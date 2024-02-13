package models

import (
	"encoding/json"

	"github.com/dgraph-io/badger"
)

type UserFriends struct {
	UserId           string   `json:"id"`
	FollowerUserIds  []string `json:"followerUserIds"`
	FollowingUserIds []string `json:"followingUserIds"`
}

type UserFriendsDB interface {
	ByUserID(userId string) (UserFriends, error)

	// Methods for altering contents
	Create(userFriends UserFriends) error
	Append(userFriends UserFriends) error
	RemoveFollower(userId, followerUserId string) error
	RemoveFollowing(userId, followerId string) error
	Update(userFriends UserFriends) error
	Delete(userFriendsId string) error
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

	followers, err := db.ByUserID(userFriends.UserId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	userFriends.FollowerUserIds = appendToSliceWithoutDuplicates(followers.FollowerUserIds, userFriends.FollowerUserIds...)
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
