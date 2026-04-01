package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func addTagRoutes(rg *gin.RouterGroup, db *sql.DB) {
	tags := rg.Group("/tags")

	// Get all tags
	tags.GET("/", func(c *gin.Context) {
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

}
