package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/seanhalberthal/webmart/docs"
	"github.com/seanhalberthal/webmart/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2" // http-swagger middleware
	"go.uber.org/zap"
	"net/http"
	"time"
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

type config struct {
	addr   string
	apiURL string
	db     dbConfig
	env    string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // React frontend
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}))

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

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsURL)))

		r.Route("/products", func(r chi.Router) {
			r.Post("/", app.createProductHandler)
			r.Get("/", app.getAllProductsHandler)

			r.Route("/{productID}", func(r chi.Router) {
				r.Get("/", app.getProductHandler)
				r.Delete("/", app.deleteProductHandler)
				r.Patch("/", app.updateProductHandler)

			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.createUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Get("/", app.getUserHandler)
				//r.Delete("/", app.deleteUserHandler)
				//r.Patch("/", app.updateUserHandler)

			})
		})

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
		})
	})

	return mux
}

func (app *application) serve(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("listening on", "addr", app.config.addr, "env", app.config.env)

	return srv.ListenAndServe()
}
