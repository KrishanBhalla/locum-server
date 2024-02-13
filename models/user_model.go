package models

import (
	"encoding/json"
	"time"

	"github.com/dgraph-io/badger"
)

type User struct {
	Id            string    `json:"id"`
	Email         string    `json:"email"`
	FullName      string    `json:"fullName"`
	LastLoginTime time.Time `json:"lastLoginTime"`
	CreationTime  time.Time `json:"creationTime"`
}

type UserDB interface {
	ByID(userId string) (User, error)
	Query(queryString string) ([]User, error)

	// Methods for altering contents
	Create(user User) error
	Update(user User) error
	Delete(userId string) error
	DbCloser
}

// Define userDB and ensure it implements UserDB
var _ UserDB = &userDB{}

type userDB struct {
	db *badger.DB
}

func NewUserDB(db *badger.DB) UserDB {
	return &userDB{db}
}

// ByID implements UserDB.
func (db *userDB) ByID(userId string) (User, error) {

	var user User
	var data = make([]byte, 0)
	data, err := lookupByKey(db.db, userId, data)
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(data, &user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// Create implements UserDB.
func (db *userDB) Create(user User) error {
	return db.Update(user)
}

// Delete implements UserDB.
func (db *userDB) Delete(userId string) error {
	user, err := db.ByID(userId)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	if err == badger.ErrKeyNotFound {
		return nil
	}

	err = db.db.Update(func(txn *badger.Txn) error {

		userBytes, err := json.Marshal(user)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(user.Id), userBytes)
		return err
	})
	return err
}

// Update implements UserDB.
func (db *userDB) Update(user User) error {
	oldUser, err := db.ByID(user.Id)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	if err == nil {
		user.CreationTime = oldUser.CreationTime
		user.LastLoginTime = time.Now()
	} else {
		user.CreationTime = time.Now()
		user.LastLoginTime = user.CreationTime
	}

	err = db.db.Update(func(txn *badger.Txn) error {

		userBytes, err := json.Marshal(user)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(user.Id), userBytes)
		return err
	})
	return err
}

// Query implements UserDB.
func (db *userDB) Query(queryString string) ([]User, error) {
	users := make([]User, 0, 0)
	var data = make([]byte, 0)
	var queryBytes []byte = []byte(queryString)
	data, err := lookupByPrefix(db.db, queryBytes, data)
	if err != nil {
		return users, err
	}
	err = json.Unmarshal(data, &users)
	if err != nil {
		return users, err
	}
	return users, nil
}

func (db *userDB) CloseDB() error {
	return db.db.Close()
}
