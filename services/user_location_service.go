package services

import (
	"time"

	"github.com/KrishanBhalla/locum-server/models"
	badger "github.com/dgraph-io/badger/v4"
)

var _ models.UserLocationDB = &userLocationService{}

var localCache map[string]time.Time = make(map[string]time.Time)

// UserLocationService is a set of methods used to manipulate
// and work with the User Location model
type UserLocationService interface {
	models.UserLocationDB
}

type userLocationService struct {
	models.UserLocationDB
}

// NewContentService initialises a ContentService object with an open connection
// to the db.
func NewUserLocationService(db *badger.DB) UserLocationService {
	return &userLocationService{models.NewUserLocationDB(db)}
}
