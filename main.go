package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/KrishanBhalla/locum-server/routes"
	"github.com/KrishanBhalla/locum-server/services"
	"github.com/KrishanBhalla/locum-server/services/websocket_service"
	"github.com/go-chi/chi"
)

const (
	port = ":8080"
)

func isProd() bool {
	return false
}

func main() {

	r := chi.NewRouter()

	services, err := services.NewServices(
		services.WithUser(),
		services.WithUserFriends(),
		services.WithUserLocation(),
	)
	defer services.Close()
	must(err)

	setupRoutes(r, *services)

	// Listen And Serve
	fmt.Println("Listening on", port)
	// Apply this to every request
	err = http.ListenAndServe(port, r) //csrfMw(r))
	must(err)
}

func setupRoutes(r *chi.Mux, services services.Services) {
	// Websocket
	r.MethodFunc(http.MethodGet, "/updateLocationWs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, &services)
	}))

	// Users
	r.MethodFunc(http.MethodPost, "/login", func(w http.ResponseWriter, r *http.Request) {
		routes.SignupOrLogin(services, w, r)
	})

	r.MethodFunc(http.MethodPost, "/friends", func(w http.ResponseWriter, r *http.Request) {
		routes.FindFriends(services, w, r)
	})

}

func serveWs(w http.ResponseWriter, r *http.Request, services *services.Services) {

	log.Println("WebSocket Endpoint Hit")
	conn, err := websocket_service.Upgrade(w, r)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%+V\n", err)
	}
	services.UserLocation.SubscribeToLocationUpdates(conn)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
