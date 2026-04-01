package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func addNoteRoutes(rg *gin.RouterGroup, db *sql.DB) {
	notes := rg.Group("/notes")

	// Get all notes
	notes.GET("/", func(c *gin.Context) {
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
}
