package services

import (
	"github.com/KrishanBhalla/locum-server/models"
	badger "github.com/dgraph-io/badger/v4"
)

var _ models.UserTokenDB = &userTokenService{}

// UserService is a set of methods used to manipulate
// and work with the user model
type UserTokenService interface {
	models.UserTokenDB
}

type userTokenService struct {
	models.UserTokenDB
}

// NewContentService initialises a ContentService object with an open connection
// to the db.
func NewUserTokenService(db *badger.DB) UserTokenService {
	return models.NewUserTokenDB(db)
}
