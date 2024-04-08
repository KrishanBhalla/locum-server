package models

import (
	"testing"
	"time"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TEST_USER_ID = "TEST"
const TEST_USER_NAME = "TEST NAME"

var TEST_USER = User{
	Id:            TEST_USER_ID,
	Email:         "Test",
	FullName:      TEST_USER_NAME,
	LastLoginTime: time.Now(),
	CreationTime:  time.Now(),
}

func TestCanCreateUserEntry(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userDb := NewUserDB(db)
		require.NoError(t, userDb.Create(TEST_USER), "Failed to create test user")
	})
}

func TestCanGetUserByID(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userDb := NewUserDB(db)
		_, err := userDb.ByID(TEST_USER_ID)
		require.ErrorContains(t, err, badger.ErrKeyNotFound.Error(), "Failed to return ErrKeyNotFound when calling ByID")

		err = userDb.Create(TEST_USER)

		user, err := userDb.ByID(TEST_USER_ID)
		require.NoError(t, err, "Failed to find test user")
		assert.EqualExportedValues(t, TEST_USER, user)
	})
}

func TestCanDeleteUser(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userDb := NewUserDB(db)

		userDb.Create(TEST_USER)

		err := userDb.Delete(TEST_USER_ID)
		require.NoError(t, err, "Failed to delete test user")
	})
}

func TestCanQueryForUser(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userDb := NewUserDB(db)
		users, err := userDb.Query(TEST_USER_NAME)
		require.NoError(t, err, "Unexpected error when querying empty db")
		assert.Empty(t, users, "Users found in empty DB")

		userDb.Create(TEST_USER)

		users, err = userDb.Query(TEST_USER_NAME)
		require.NoError(t, err, "Unexpected error when querying non-empty db")
		assert.NotEmpty(t, users, "Users not found in DB")
	})
}

func TestCanUpdateUserEmail(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userDb := NewUserDB(db)
		userDb.Create(TEST_USER)

		modifiedTestUser := TEST_USER
		modifiedTestUser.Email = "New Email For Unit Test"
		err := userDb.Update(modifiedTestUser)
		require.NoError(t, err, "Unexpected error when editing user")

		user, err := userDb.ByID(TEST_USER_ID)
		require.NoError(t, err, "Failed to find test user")
		assert.EqualExportedValues(t, modifiedTestUser, user)
	})
}
