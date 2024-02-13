package services

import "github.com/KrishanBhalla/locum-server/models"

// Services contains the services which need db connection
type Services struct {
	User         UserService
	UserFriends  UserFriendsService
	UserLocation UserLocationService
}

// NewServices initialises all services with a single db connection
func NewServices(cfgs ...ServicesConfig) (*Services, error) {

	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

// Close closes the database connections.
func (s *Services) Close() error {
	closers := []models.DbCloser{s.User, s.UserFriends, s.UserLocation}
	for _, c := range closers {
		err := c.CloseDB()
		if err != nil {
			return err
		}
	}
	return nil
}
