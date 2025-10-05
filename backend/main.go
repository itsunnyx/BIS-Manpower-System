package main

import (
	httpTransport "manpower/transport/http"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgres://user:password@localhost:5432/manpower?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	srv := httpTransport.NewServer(db)
	log.Println("Server started at :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
