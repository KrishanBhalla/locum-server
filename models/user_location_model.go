package models

import (
	"encoding/json"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

type GeoTime struct {
	Latitude  float32   `json:"latitude"`
	Longitude float32   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}

type UserLocation struct {
	UserId   string    `json:"userId"`
	GeoTimes []GeoTime `json:"geoTime"`
}

type UserLocationDB interface {
	ByUserID(userId string) (UserLocation, error)
	LatestGeoTimeByUserID(userId string) (GeoTime, error)

	// Methods for altering UserLocations
	Create(userLocation UserLocation) error
	Append(userLocation UserLocation) error
	Update(userLocation UserLocation) error
	Delete(userId string) error
	DbCloser
}

// Define userLocationDB and ensure it implements UserLocationDB
var _ UserLocationDB = &userLocationDB{}

func NewUserLocationDB(db *badger.DB) UserLocationDB {
	return &userLocationDB{db}
}

type userLocationDB struct {
	db *badger.DB
}

// ByUserID implements UserLocationDB.
func (db *userLocationDB) ByUserID(userId string) (UserLocation, error) {

	var userLocation = UserLocation{}
	var data = make([]byte, 0)

	data, err := lookupByKey(db.db, userId, data)
	if err != nil {
		return UserLocation{}, err
	}
	err = json.Unmarshal(data, &userLocation)
	if err != nil {
		return UserLocation{}, err
	}
	return userLocation, nil
}

// LatestByUserID implements UserLocationDB.
func (db *userLocationDB) LatestGeoTimeByUserID(userId string) (GeoTime, error) {
	userLocation, err := db.ByUserID(userId)
	if err != nil || len(userLocation.GeoTimes) == 0 {
		return GeoTime{}, err
	}
	return userLocation.GeoTimes[len(userLocation.GeoTimes)-1], nil
}

// Create will create the provided UserLocation and backfill data
func (db *userLocationDB) Create(userLocation UserLocation) error {
	return db.Update(userLocation)
}

func (db *userLocationDB) Append(userLocation UserLocation) error {

	existingUserLocation, err := db.ByUserID(userLocation.UserId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	} else if err == nil {
		userLocation.GeoTimes = append(existingUserLocation.GeoTimes, userLocation.GeoTimes...)
	} // key not found vacuosly works

	return db.Update(userLocation)
}

func (db *userLocationDB) Update(userLocation UserLocation) error {

	err := db.db.Update(func(txn *badger.Txn) error {
		userLocationBytes, err := json.Marshal(userLocation)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(userLocation.UserId), userLocationBytes)
		return err
	})
	return err
}

func (db *userLocationDB) Delete(userId string) error {
	err := db.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(userId))
		return err
	})
	return err
}

func (db *userLocationDB) CloseDB() error {
	return db.db.Close()
}
