package main

import (
	"github.com/bssmnt/webmart/internal/env"
	"log"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	app := &application{
		config: cfg,
	}

	mux := app.routes()

	log.Fatal(app.serve(mux))
}
