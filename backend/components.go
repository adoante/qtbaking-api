package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func addComponentRoutes(rg *gin.RouterGroup, db *sql.DB) {
	components := rg.Group("/components")

	// Get all components
	components.GET("/", func(c *gin.Context) {
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
}
