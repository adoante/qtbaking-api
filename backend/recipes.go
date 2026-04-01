package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func addRecipeRoutes(rg *gin.RouterGroup, db *sql.DB) {
	recipes := rg.Group("/recipes")

	// Get recipe by id
	recipes.GET("/", func(c *gin.Context) {
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

	// Get recipe by id
	recipes.GET("/:id", func(c *gin.Context) {
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

}
