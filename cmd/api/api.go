package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/seanhalberthal/webmart/internal/store"
	"log"
	"net/http"
	"time"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	mux.Use(middleware.Timeout(60 * time.Second))

	mux.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/products", func(r chi.Router) {
			r.Post("/", app.createProductHandler)

			r.Route("/{productID}", func(r chi.Router) {
				r.Get("/", app.getProductHandler)
				r.Delete("/", app.deleteProductHandler)
				r.Patch("/", app.updateProductHandler)

			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.createUserHandler)
		})
	})

	return mux
}

func (app *application) serve(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("listening on %s", app.config.addr)

	return srv.ListenAndServe()
}
