package main

import (
	"database/sql"
	"github.com/seanhalberthal/webmart/internal/db"
	"github.com/seanhalberthal/webmart/internal/env"
	"github.com/seanhalberthal/webmart/internal/store"
	"log"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgresql://postgres:password@localhost/webmart?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	s := store.NewStorage(conn)
	db.Seed(s)
}
