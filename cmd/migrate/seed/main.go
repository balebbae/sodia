package main

import (
	"log"

	"github.com/balebbae/sodia/internal/db"
	"github.com/balebbae/sodia/internal/env"
	"github.com/balebbae/sodia/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/socialnetwork?sslmode=disable")
	log.Println("Using DB Address:", addr) // Debugging line

	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store)
}