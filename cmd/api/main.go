package main

import (
	"database/sql"
	"fmt"
	"github.com/seanhalberthal/webmart/internal/db"
	"github.com/seanhalberthal/webmart/internal/env"
	"github.com/seanhalberthal/webmart/internal/store"
	"log"
)

const version = "0.0.1"

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgresql://postgres:password@localhost/webmart?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	database, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(database)

	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(database)
	fmt.Println("database connection established")

	storage := store.NewStorage(database)

	app := &application{
		config: cfg,
		store:  storage,
	}

	mux := app.routes()

	log.Fatal(app.serve(mux))
}
