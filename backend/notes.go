package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAllNotes(db *sql.DB) ([]Note, error) {
	rows, err := db.Query(`
		SELECT id, recipe_id, note
		FROM notes
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note

	for rows.Next() {
		var note Note

		err := rows.Scan(
			&note.ID,
			&note.RecipeId,
			&note.Note,
		)

		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func getNotesByRecipeId(db *sql.DB, recipeId string) ([]Note, error) {
	rows, err := db.Query(`
		SELECT id, recipe_id, note
		FROM notes
		WHERE recipe_id = $1
	`, recipeId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note

	for rows.Next() {
		var note Note

		err := rows.Scan(
			&note.ID,
			&note.RecipeId,
			&note.Note,
		)

		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func addNoteRoutes(rg *gin.RouterGroup) {
	notes := rg.Group("/notes")

	// Get all notes
	notes.GET("/", func(c *gin.Context) {

		db := c.MustGet("db").(*sql.DB)

		notes, err := getAllNotes(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, notes)
	})
}
