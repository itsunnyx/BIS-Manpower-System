package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"manpower/transport/http/handlers"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var db *sql.DB

func initDB() *sql.DB {
	// ตัวอย่าง connection string
	dsn := "host=localhost port=5432 user=admin password=1234 dbname=hr sslmode=disable"

	database, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// ตรวจสอบการเชื่อมต่อ
	err = database.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Connected to database successfully!")
	return database
}

func main() {
	db = initDB()
	defer db.Close()
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// routes จริง
	api := r.Group("/api")
	{
		api.GET("/requests", handlers.GetManpowerRequests(db))    // ดึงรายการคำร้อง
		api.POST("/requests", handlers.CreateManpowerRequest(db)) // สร้างคำร้องใหม่
	}

	// start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server started at :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
