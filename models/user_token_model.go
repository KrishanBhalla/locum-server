package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

const SECRET_PEPPER = "locum-secret-pepper"

type UserToken struct {
	UserId       string    `json:"userId"`
	Token        string    `json:"token"`
	CreationTime time.Time `json:"creationTime"`
}

func NewUserToken(userId string) UserToken {

	timeNow := time.Now().UTC()

	rawTokenBytes := append([]byte(SECRET_PEPPER), userId...)
	rawTokenBytes = append(rawTokenBytes, []byte(timeNow.String())...)
	h := sha256.New()
	h.Write([]byte(rawTokenBytes))
	token := hex.EncodeToString(h.Sum(nil))

	return UserToken{UserId: userId, Token: token, CreationTime: timeNow}
}

type UserTokenDB interface {
	ByToken(token string) (UserToken, error)
	// Methods for altering contents
	Create(userToken UserToken) error
	Update(userToken UserToken) error
	Delete(userId string) error
	DbCloser
}

// Define userTokenDB and ensure it implements UserTokenModel
var _ UserTokenDB = &userTokenDB{}

type userTokenDB struct {
	db *badger.DB
}

func NewUserTokenDB(db *badger.DB) UserTokenDB {
	return &userTokenDB{db}
}

// ByID implements UserTokenModel.
func (db *userTokenDB) ByToken(token string) (UserToken, error) {

	var userToken UserToken
	var data = make([]byte, 0)
	data, err := lookupByKey(db.db, token, data)
	if err != nil {
		return UserToken{}, err
	}
	err = json.Unmarshal(data, &userToken)
	if err != nil {
		return UserToken{}, err
	}
	return userToken, nil
}

// Create implements UserTokenModel.
func (db *userTokenDB) Create(userToken UserToken) error {

	// First clear all existing tokens
	tokens, err := db.byUserId(userToken.UserId)
	if err != nil {
		return err
	}
	for _, token := range tokens {
		err = db.Delete(token.Token)
		if err != nil {
			return err
		}
	}

	err = db.Update(userToken)
	if err != nil {
		return err
	}
	return nil
}

// Delete implements UserTokenModel.
func (db *userTokenDB) Delete(token string) error {
	userToken, err := db.ByToken(token)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	if err == badger.ErrKeyNotFound {
		return nil
	}

	err = db.db.Update(func(txn *badger.Txn) error {
		err = txn.Delete([]byte(userToken.Token))
		return err
	})
	return err
}

// Update implements UserTokenModel.
func (db *userTokenDB) Update(userToken UserToken) error {
	oldUserToken, err := db.ByToken(userToken.Token)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	if err == nil {
		userToken.CreationTime = oldUserToken.CreationTime
	} else {
		userToken.CreationTime = time.Now()
	}

	err = db.db.Update(func(txn *badger.Txn) error {

		userTokenBytes, err := json.Marshal(userToken)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(userToken.Token), userTokenBytes)
		return err
	})
	return err
}

func (db *userTokenDB) byUserId(userId string) ([]UserToken, error) {

	var data = make([]UserToken, 0)
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				u := UserToken{}
				err := json.Unmarshal(v, &u)
				if err != nil {
					return err
				}
				if strings.Contains(u.UserId, userId) {
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

func (db *userTokenDB) CloseDB() error {
	return db.db.Close()
}
