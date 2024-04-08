package models

import (
	"testing"
	"time"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TEST_FRIEND_ID = "Test 1"
const TEST_FRIEND_REQUEST_ID = "Test 2"

var TEST_USER_FRIENDS = UserFriends{
	UserId:         TEST_USER_ID,
	FriendIds:      []string{},
	FriendRequests: []FriendRequest{},
}
var TEST_FRIEND = UserFriends{
	UserId:         TEST_FRIEND_ID,
	FriendIds:      []string{TEST_USER_ID},
	FriendRequests: []FriendRequest{},
}

var TEST_FRIEND_REQUEST = FriendRequest{
	UserId:    TEST_FRIEND_REQUEST_ID,
	Timestamp: time.Now(),
}

func TestCanCreateUserFriendsEntry(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userFriendsDb := NewUserFriendsDB(db)
		require.NoError(t, userFriendsDb.Create(TEST_USER_FRIENDS), "Failed to create test friends")
	})
}

func TestCanGetUserFriendsByUserID(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userFriendsDb := NewUserFriendsDB(db)
		_, err := userFriendsDb.ByUserID(TEST_USER_ID)
		require.ErrorContains(t, err, badger.ErrKeyNotFound.Error(), "Failed to return ErrKeyNotFound when calling ByUserID")

		err = userFriendsDb.Create(TEST_USER_FRIENDS)

		userFriends, err := userFriendsDb.ByUserID(TEST_USER_ID)
		require.NoError(t, err, "Failed to find test user friends")
		assert.EqualExportedValues(t, TEST_USER_FRIENDS, userFriends)
	})
}

func TestCanDeleteUserFriends(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userFriendsDb := NewUserFriendsDB(db)

		userFriendsDb.Create(TEST_USER_FRIENDS)

		err := userFriendsDb.Delete(TEST_USER_ID)
		require.NoError(t, err, "Failed to delete test user friends")
	})
}

func TestCanAddFriends(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userFriendsDb := NewUserFriendsDB(db)
		userFriendsDb.Create(TEST_USER_FRIENDS)

		userFriendsDb.Create(TEST_FRIEND)
		err := userFriendsDb.AddFriend(TEST_USER_ID, TEST_FRIEND_ID)
		require.NoError(t, err, "Unexpected error when adding friend")

		userFriends, err := userFriendsDb.ByUserID(TEST_USER_ID)

		assert.Greater(t, len(userFriends.FriendIds), len(TEST_USER_FRIENDS.FriendIds), "No new friend added")
		assert.Contains(t, userFriends.FriendIds, TEST_FRIEND_ID, "New friend missing")
	})
}

func TestCanRemoveFriends(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userFriendsDb := NewUserFriendsDB(db)
		userFriendsDb.Create(TEST_USER_FRIENDS)
		userFriendsDb.Create(TEST_FRIEND)
		userFriendsDb.AddFriend(TEST_USER_ID, TEST_FRIEND_ID)

		err := userFriendsDb.RemoveFriend(TEST_USER_ID, TEST_FRIEND_ID)
		require.NoError(t, err, "Unexpected error when removing friend")

		userFriends, err := userFriendsDb.ByUserID(TEST_USER_ID)

		assert.Equal(t, len(userFriends.FriendIds), len(TEST_USER_FRIENDS.FriendIds), "No friend removed")
		assert.NotContains(t, userFriends.FriendIds, TEST_FRIEND_ID, "Old friend not removed")
	})
}

func TestCanAddFriendRequest(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userFriendsDb := NewUserFriendsDB(db)
		userFriendsDb.Create(TEST_USER_FRIENDS)

		err := userFriendsDb.AddFriendRequest(TEST_USER_ID, TEST_FRIEND_REQUEST)
		require.NoError(t, err, "Unexpected error when adding friend request")

		userFriends, err := userFriendsDb.ByUserID(TEST_USER_ID)

		assert.Greater(t, len(userFriends.FriendRequests), len(TEST_USER_FRIENDS.FriendRequests), "No new friendRequest added")
		assert.EqualExportedValues(t, userFriends.FriendRequests[0], TEST_FRIEND_REQUEST, "New friendRequest missing")

		err = userFriendsDb.AddFriendRequest(TEST_USER_ID, TEST_FRIEND_REQUEST)
		require.NoError(t, err, "Unexpected error when adding repeat friend request")

		userFriendsWithDoubleRequest, err := userFriendsDb.ByUserID(TEST_USER_ID)

		assert.Equal(t, len(userFriendsWithDoubleRequest.FriendRequests), len(userFriends.FriendRequests), "New friendRequest erroneously added")

	})
}

func TestCanRemoveFriendRequest(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userFriendsDb := NewUserFriendsDB(db)
		userFriendsDb.Create(TEST_USER_FRIENDS)
		userFriendsDb.AddFriendRequest(TEST_USER_ID, TEST_FRIEND_REQUEST)

		err := userFriendsDb.RemoveFriendRequest(TEST_USER_ID, TEST_FRIEND_REQUEST_ID)

		require.NoError(t, err, "Unexpected error when removing friend request")

		userFriends, _ := userFriendsDb.ByUserID(TEST_USER_ID)

		assert.Equal(t, len(userFriends.FriendRequests), len(TEST_USER_FRIENDS.FriendRequests), "Friend Request not removed")
		assert.NotContains(t, userFriends.FriendRequests, TEST_FRIEND_REQUEST, "FriendRequest remains")
	})
}
