package main

import (
	"log"
	"manpower/internal/db"
	"manpower/internal/server"
)

func main() {
	// เริ่มต้น DB
	conn := db.InitDB()
	defer conn.Close()

	// สร้าง server
	srv := server.NewServer(conn)

	// รัน server
	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}