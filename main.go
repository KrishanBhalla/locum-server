package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/KrishanBhalla/locum-server/api"
	"github.com/KrishanBhalla/locum-server/api/spec"
	"github.com/KrishanBhalla/locum-server/middleware"
	"github.com/KrishanBhalla/locum-server/services"
	"github.com/KrishanBhalla/locum-server/services/websocket_service"
	"github.com/go-chi/chi"
	chiMw "github.com/go-chi/chi/middleware"
)

const (
	port = ":8080"
)

func isProd() bool {
	return false
}

func main() {

	services, err := services.NewServices(
		services.WithUser(),
		services.WithUserFriends(),
		services.WithUserLocation(),
	)
	defer services.Close()
	must(err)

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(chiMw.RequestID)
	r.Use(chiMw.RealIP)
	r.Use(chiMw.Logger)
	r.Use(chiMw.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(chiMw.Timeout(60 * time.Second))
	r.Use(middleware.AddServices(services))

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

	r.MethodFunc(http.MethodGet, "/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Successful test")
	}))
	r.Mount("/api", spec.Handler(spec.NewStrictHandler(&api.ServerImpl{}, []spec.StrictMiddlewareFunc{})))

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
