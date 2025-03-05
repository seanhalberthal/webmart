package main

import (
	"database/sql"
	"fmt"
	"github.com/bssmnt/webmart/internal/env"
	"github.com/bssmnt/webmart/internal/store"
	"log"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgresql://postgres:password@localhost/webmart?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	database, err := sql.Open("postgres", cfg.db.addr)
	if err != nil {
		log.Fatal(err)
	}

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
