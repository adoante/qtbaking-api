package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

func main() {
	// Load Enviorment variables
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbHost := os.Getenv("POSTGRES_HOST")

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	const maxAttempts = 10
	for i := 1; i <= maxAttempts; i++ {
		err = db.Ping()
		if err == nil {
			fmt.Println("Connected!")
			break
		}

		log.Printf("Database not ready yet (attempt %d/%d): %v", i, maxAttempts, err)
		time.Sleep(2 * time.Second)
	}

	fmt.Println("Connected!")

	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Get Routes
	v1 := r.Group("/v1")

	addVodRoutes(v1, db)
	addRecipeRoutes(v1, db)
	addComponentRoutes(v1, db)
	addIngredientsRoutes(v1, db)
	addToolsRoutes(v1, db)
	addNoteRoutes(v1, db)
	addTagRoutes(v1, db)
	addBakealongRoutes(v1, db)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
