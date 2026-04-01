package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func addToolsRoutes(rg *gin.RouterGroup, db *sql.DB) {
	notes := rg.Group("/tools")

	// Get all notes
	notes.GET("/", func(c *gin.Context) {
		rows, err := db.Query(
			`SELECT id, recipe_id, name, Optional
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

		c.JSON(http.StatusOK, notes)
	})
}
