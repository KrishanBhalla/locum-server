package models

import (
	"crypto/sha256"
	"encoding/json"
	"strings"
	"time"

	"github.com/dgraph-io/badger"
)

type User struct {
	Id             string    `json:"id"`
	Email          string    `json:"email"`
	FullName       string    `json:"fullName"`
	LastLoginTime  time.Time `json:"lastLoginTime"`
	CreationTime   time.Time `json:"creationTime"`
	ShareableToken []byte    `json:"shareableToken"`
}

type UserDB interface {
	ByID(userId string) (User, error)
	Query(queryString string) ([]User, error)

	// Methods for altering contents
	Create(user User) (string, error)
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
func (db *userDB) Create(user User) (string, error) {
	err := db.Update(user)
	if err != nil {
		return "", err
	}
	user, err = db.ByID(user.Id)
	if err != nil {
		return "", err
	}
	return user.Id, nil
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

	if user.ShareableToken == nil {
		token := sha256.New().Sum([]byte(user.Id))
		user.ShareableToken = token
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

	var data = make([]User, 0)
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				u := User{}
				err := json.Unmarshal(v, &u)
				if err != nil {
					return err
				}
				if strings.Contains(u.FullName, queryString) {
					data = append(data, u)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return data, nil
}

// QueryByToken implements UserDB.
func (db *userDB) QueryByToken(tokenBytes []byte) (User, error) {

	var data = &User{}
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		continueIteration := true
		for it.Rewind(); it.Valid() && continueIteration; it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				u := User{}
				err := json.Unmarshal(v, &u)
				if err != nil {
					return err
				}
				if string(u.ShareableToken) == string(tokenBytes) {
					data = &u
					continueIteration = false
					return nil
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return *data, err
	}
	return *data, nil
}

func (db *userDB) CloseDB() error {
	return db.db.Close()
}
