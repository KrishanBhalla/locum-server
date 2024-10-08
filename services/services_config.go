package services

import (
	badger "github.com/dgraph-io/badger/v4"
)

// ServicesConfig allows for dynamic adding of services
type ServicesConfig func(*Services) error

// withBadger initiates a badger db
func withBadger(dbPath string) (*badger.DB, error) {
	opt := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// WithUser returns a ServicesConfig object that sets a user
func WithUser() ServicesConfig {
	return func(s *Services) error {
		db, err := withBadger("badger/users")
		if err != nil {
			return err
		}
		s.User = NewUserService(db)
		return nil
	}
}

// WithUserFriends returns a ServicesConfig object that sets a user
func WithUserFriends() ServicesConfig {
	return func(s *Services) error {
		db, err := withBadger("badger/userFriends")
		if err != nil {
			return err
		}
		s.UserFriends = NewUserFriendsService(db)
		return nil
	}
}

// WithUserLocation returns a ServicesConfig object that sets a userLocation
func WithUserLocation() ServicesConfig {
	return func(s *Services) error {
		db, err := withBadger("badger/userLocation")
		if err != nil {
			return err
		}
		s.UserLocation = NewUserLocationService(db)
		return nil
	}
}

// WithUserToken returns a ServicesConfig object that sets a User token
func WithUserToken() ServicesConfig {
	return func(s *Services) error {
		db, err := withBadger("badger/userToken")
		if err != nil {
			return err
		}
		s.UserToken = NewUserTokenService(db)
		return nil
	}
}
