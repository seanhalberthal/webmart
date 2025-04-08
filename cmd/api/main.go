package main

import (
	"database/sql"
	"github.com/seanhalberthal/webmart/internal/db"
	"github.com/seanhalberthal/webmart/internal/env"
	"github.com/seanhalberthal/webmart/internal/store"
	"go.uber.org/zap"
	"log"
)

const version = "0.0.1"

//	@title			Webmart
//	@description	API for Webmart, an e-commerce solution
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {
	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgresql://postgres:password@localhost/webmart?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer func(logger *zap.SugaredLogger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)

	// Database
	database, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}(database)
	logger.Info("database connection established")

	storage := store.NewStorage(database)

	app := &application{
		config: cfg,
		store:  storage,
		logger: logger,
	}

	mux := app.routes()

	logger.Fatal(app.serve(mux))
}
