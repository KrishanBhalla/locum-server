package models

import (
	"testing"
	"time"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var TEST_USER_LOCATION = UserLocation{
	UserId: TEST_USER_ID,
	GeoTimes: []GeoTime{
		{
			Latitude:  0,
			Longitude: 0,
			Timestamp: time.Unix(0, 0),
		},
		{
			Latitude:  1,
			Longitude: 0,
			Timestamp: time.Unix(1, 0),
		},
	},
}

var TEST_GEOTIME = GeoTime{
	Latitude:  1,
	Longitude: 1,
	Timestamp: time.Unix(2, 0),
}

func TestCanCreateUserLocationEntry(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userLocationDb := NewUserLocationDB(db)
		require.NoError(t, userLocationDb.Create(TEST_USER_LOCATION), "Failed to create test user location")
	})
}

func TestCanAddGeotime(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userLocationDb := NewUserLocationDB(db)
		userLocationDb.Create(TEST_USER_LOCATION)

		newUserLocation := UserLocation{
			UserId:   TEST_USER_ID,
			GeoTimes: []GeoTime{TEST_GEOTIME},
		}
		err := userLocationDb.Append(newUserLocation)
		require.NoError(t, err, "Failed to append new user location")

		g, err := userLocationDb.LatestGeoTimeByUserID(TEST_USER_ID)
		require.NoError(t, err, "Failed to find test user location")
		assert.EqualExportedValues(t, g, TEST_GEOTIME, "Did not add expected geotime")
	})
}

func TestCanGetUserLocationByUserID(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userLocationDb := NewUserLocationDB(db)
		_, err := userLocationDb.ByUserID(TEST_USER_ID)
		require.ErrorContains(t, err, badger.ErrKeyNotFound.Error(), "Failed to return ErrKeyNotFound when calling ByUserID")

		err = userLocationDb.Create(TEST_USER_LOCATION)

		userLoc, err := userLocationDb.ByUserID(TEST_USER_ID)
		require.NoError(t, err, "Failed to find test user location")
		for i, g := range TEST_USER_LOCATION.GeoTimes {
			assert.EqualExportedValues(t, userLoc.GeoTimes[i], g)
		}
	})
}

func TestCanGetLatestUserLocationByUserID(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userLocationDb := NewUserLocationDB(db)
		userLocationDb.Create(TEST_USER_LOCATION)

		geoTime, err := userLocationDb.LatestGeoTimeByUserID(TEST_USER_ID)
		require.NoError(t, err, "Failed to find test user location")
		assert.EqualExportedValues(t, geoTime, TEST_USER_LOCATION.GeoTimes[len(TEST_USER_LOCATION.GeoTimes)-1])

	})
}

func TestCanDeleteUserLocation(t *testing.T) {

	runBadgerTest(t, func(t *testing.T, db *badger.DB) {
		userLocationDb := NewUserLocationDB(db)

		userLocationDb.Create(TEST_USER_LOCATION)

		err := userLocationDb.Delete(TEST_USER_ID)
		require.NoError(t, err, "Failed to delete test user locatino")
	})
}
