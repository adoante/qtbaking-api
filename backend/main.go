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

type Vod struct {
	ID        int
	Slug      string
	Title     string
	VideoURL  string
	CreatedAt time.Time
}

type Recipe struct {
	ID             int
	VodId          sql.NullInt64
	Thumbnail      sql.NullString
	Title          string
	TempFahrenheit sql.NullInt64
	TempCelsius    sql.NullInt64
}

type Component struct {
	ID       int
	RecipeId int
	Name     string
}

type Ingredient struct {
	ID             int
	ComponentId    int
	Name           string
	Quantity       float64
	Unit           string
	MetricQuantity sql.NullFloat64
	MetricUnit     sql.NullString
	Optional       bool
	Notes          string
}

type Tool struct {
	ID       int
	RecipeId int
	Name     string
	Optional bool
}

type Note struct {
	ID       int
	RecipeId int
	Note     string
}

type Tag struct {
	ID       int
	RecipeId int
	Tag      string
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

	// Get all vods
	// TODO: addd pagination
	r.GET("/vods", func(c *gin.Context) {
		rows, err := db.Query(
			`SELECT id, slug, title, video_url, created_at
			FROM vods
		`)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		// Return JSON response
		var vods []Vod

		for rows.Next() {
			var vod Vod
			err := rows.Scan(
				&vod.ID,
				&vod.Slug,
				&vod.Title,
				&vod.VideoURL,
				&vod.CreatedAt,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			vods = append(vods, vod)
		}

		c.JSON(http.StatusOK, vods)

	})

	// get vod by slug
	r.GET("/vods/:slug", func(c *gin.Context) {
		id := c.Param("slug")

		rows, err := db.Query(
			`SELECT id, slug, title, video_url, created_at
			FROM vods
			WHERE slug = $1
		`, id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		// Return JSON response
		var vods []Vod

		for rows.Next() {
			var vod Vod
			err := rows.Scan(
				&vod.ID,
				&vod.Slug,
				&vod.Title,
				&vod.VideoURL,
				&vod.CreatedAt,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			vods = append(vods, vod)
		}

		c.JSON(http.StatusOK, vods)

	})

	// Get recipe by id
	r.GET("/recipes/:id", func(c *gin.Context) {
		id := c.Param("id")

		rows, err := db.Query(
			`SELECT id, vod_id, title, thumbnail, temp_fahrenheit, temp_celsius
			FROM recipes
			WHERE id = $1
		`, id)

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
				&recipe.VodId,
				&recipe.Title,
				&recipe.Thumbnail,
				&recipe.TempFahrenheit,
				&recipe.TempCelsius,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			recipes = append(recipes, recipe)
		}

		c.JSON(http.StatusOK, recipes)

	})

	// Get all recipes
	r.GET("/recipes", func(c *gin.Context) {
		rows, err := db.Query(
			`SELECT id, vod_id, title, thumbnail, temp_fahrenheit, temp_celsius
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
				&recipe.VodId,
				&recipe.Title,
				&recipe.Thumbnail,
				&recipe.TempFahrenheit,
				&recipe.TempCelsius,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			recipes = append(recipes, recipe)
		}

		c.JSON(http.StatusOK, recipes)

	})

	// Get all components
	r.GET("/components", func(c *gin.Context) {
		rows, err := db.Query(
			`SELECT id, recipe_id, name
			FROM components
		`)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		// Return JSON response
		var components []Component

		for rows.Next() {
			var component Component
			err := rows.Scan(
				&component.ID,
				&component.RecipeId,
				&component.Name,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			components = append(components, component)
		}

		c.JSON(http.StatusOK, components)

	})

	// Get all ingredients
	r.GET("/ingredients", func(c *gin.Context) {
		rows, err := db.Query(
			`SELECT id, component_id, name, quantity, unit, metric_quantity, metric_unit, optional, notes
			FROM ingredients
		`)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		// Return JSON response
		var ingredients []Ingredient

		for rows.Next() {
			var ingredient Ingredient
			err := rows.Scan(
				&ingredient.ID,
				&ingredient.ComponentId,
				&ingredient.Name,
				&ingredient.Quantity,
				&ingredient.Unit,
				&ingredient.MetricQuantity,
				&ingredient.MetricUnit,
				&ingredient.Optional,
				&ingredient.Notes,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			ingredients = append(ingredients, ingredient)
		}

		c.JSON(http.StatusOK, ingredients)

	})

	// Get all tools
	r.GET("/tools", func(c *gin.Context) {
		rows, err := db.Query(
			`SELECT id, recipe_id, name, optional
			FROM tools
		`)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		// Return JSON response
		var tools []Tool

		for rows.Next() {
			var tool Tool
			err := rows.Scan(
				&tool.ID,
				&tool.RecipeId,
				&tool.Name,
				&tool.Optional,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			tools = append(tools, tool)
		}

		c.JSON(http.StatusOK, tools)

	})

	// Get all notes
	r.GET("/notes", func(c *gin.Context) {
		rows, err := db.Query(
			`SELECT id, recipe_id, note
			FROM notes
		`)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		// Return JSON response
		var notes []Note

		for rows.Next() {
			var note Note
			err := rows.Scan(
				&note.ID,
				&note.RecipeId,
				&note.Note,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			notes = append(notes, note)
		}

		c.JSON(http.StatusOK, notes)

	})

	// Get all tags
	r.GET("/tags", func(c *gin.Context) {
		rows, err := db.Query(
			`SELECT id, recipe_id, tag
			FROM tags
		`)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		// Return JSON response
		var tags []Tag

		for rows.Next() {
			var tag Tag
			err := rows.Scan(
				&tag.ID,
				&tag.RecipeId,
				&tag.Tag,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			tags = append(tags, tag)
		}

		c.JSON(http.StatusOK, tags)

	})
	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
