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

// https://gin-gonic.com/en/docs/server-config/database/#middleware
func DatabaseMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

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

	// Retry connection
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
	r.Use(DatabaseMiddleware(db))

	// Add Routes
	v1 := r.Group("/")

	addVodRoutes(v1)
	addRecipeRoutes(v1)
	addComponentRoutes(v1)
	addIngredientRoutes(v1)
	addToolRoutes(v1)
	addNoteRoutes(v1)
	addTagRoutes(v1)
	addBakealongRoutes(v1)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
