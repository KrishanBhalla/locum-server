package services

import (
	"github.com/KrishanBhalla/locum-server/models"
	"github.com/dgraph-io/badger"
)

var _ models.UserDB = &userService{}

// UserService is a set of methods used to manipulate
// and work with the user model
type UserService interface {
	models.UserDB
}

type userService struct {
	models.UserDB
}

// NewContentService initialises a ContentService object with an open connection
// to the db.
func NewUserService(db *badger.DB) UserService {
	return models.NewUserDB(db)
}
