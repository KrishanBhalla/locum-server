package websocket_service

type LocationUpdateMessage struct {
	UserID    string `json:"userId"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Timestamp string `json:"timestamp"`
}
