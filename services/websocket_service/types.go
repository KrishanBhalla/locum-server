package websocket_service

type LocationUpdateMessage struct {
	UserID    string  `json:"userId"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}
