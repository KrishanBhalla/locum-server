package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/KrishanBhalla/locum-server/api"
	"github.com/KrishanBhalla/locum-server/api/spec"
	"github.com/KrishanBhalla/locum-server/middleware"
	"github.com/KrishanBhalla/locum-server/services"
	"github.com/go-chi/chi/v5"
	chiMw "github.com/go-chi/chi/v5/middleware"
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
		services.WithUserToken(),
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

	setupRoutes(r)

	// Listen And Serve
	fmt.Println("Listening on", port)
	// Apply this to every request
	err = http.ListenAndServe(port, r) //csrfMw(r))
	must(err)
}

func setupRoutes(r *chi.Mux) {

	r.MethodFunc(http.MethodGet, "/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Successful test")
	}))
	r.Mount("/api", spec.Handler(spec.NewStrictHandler(&api.ServerImpl{}, []spec.StrictMiddlewareFunc{})))

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
