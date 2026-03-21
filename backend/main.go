package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type Recipe struct {
	ID             int
	Slug           string
	Title          string
	Thumbnail      string
	TempFahrenheit int
	TempCelsius    int
	VideoURL       string
	CreatedAt      time.Time
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

	// Get all recipes
	r.GET("/recipes", func(c *gin.Context) {
		rows, err := db.Query(
			`SELECT id, slug, title, thumbnail, temp_fahrenheit, temp_celsius, video_url
			FROM recipes
		`)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		// Return JSON response
		var recipes []Recipe

		for rows.Next() {
			var recipe Recipe
			err := rows.Scan(
				&recipe.ID,
				&recipe.Slug,
				&recipe.Title,
				&recipe.Thumbnail,
				&recipe.TempFahrenheit,
				&recipe.TempCelsius,
				&recipe.VideoURL,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			recipes = append(recipes, recipe)
		}

		c.JSON(http.StatusOK, recipes)

	})

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
