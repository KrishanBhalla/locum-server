package services

import (
	"github.com/KrishanBhalla/locum-server/models"
	badger "github.com/dgraph-io/badger/v4"
)

var _ models.UserFriendsDB = &userFriendsService{}

// UserFriendsService is a set of methods used to manipulate
// and work with the User Friends model
type UserFriendsService interface {
	models.UserFriendsDB
}

type userFriendsService struct {
	models.UserFriendsDB
}

// NewContentService initialises a ContentService object with an open connection
// to the db.
func NewUserFriendsService(db *badger.DB) UserFriendsService {
	return models.NewUserFriendsDB(db)
}
