package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/KrishanBhalla/locum-server/models"
	"github.com/KrishanBhalla/locum-server/services/websocket_service"
	"github.com/dgraph-io/badger"
	"github.com/gorilla/websocket"
)

var _ models.UserLocationDB = &userLocationService{}

// UserLocationService is a set of methods used to manipulate
// and work with the User Location model
type UserLocationService interface {
	models.UserLocationDB
	SubscribeToLocationUpdates(wsConn *websocket.Conn)
}

type userLocationService struct {
	models.UserLocationDB
}

// NewContentService initialises a ContentService object with an open connection
// to the db.
func NewUserLocationService(db *badger.DB) UserLocationService {
	return &userLocationService{models.NewUserLocationDB(db)}
}

func (s *userLocationService) SubscribeToLocationUpdates(wsConn *websocket.Conn) {

	defer func() {
		wsConn.Close()
	}()
	logger := log.Default()
	for {

		_, p, err := wsConn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var message websocket_service.LocationUpdateMessage
		err = json.Unmarshal(p, &message)
		if err != nil {
			logger.Println("Websocket (user_location_service.go):", err)
			return
		}
		// We expect an ISO string
		messageTime, err := time.Parse("2011-10-05T14:48:00.000Z", message.Timestamp)
		if err != nil {
			messageTime = time.Now()
		}
		geoTime := models.GeoTime{
			Latitude:  message.Latitude,
			Longitude: message.Longitude,
			Timestamp: messageTime,
		}
		s.Append(models.UserLocation{UserId: message.UserID, GeoTimes: []models.GeoTime{geoTime}})
	}
}
