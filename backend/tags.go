package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAllTags(db *sql.DB) ([]Tag, error) {
	rows, err := db.Query(`
		SELECT id, recipe_id, tag
		FROM tags
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag

	for rows.Next() {
		var tag Tag

		err := rows.Scan(
			&tag.ID,
			&tag.RecipeId,
			&tag.Tag,
		)

		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func getTagsByRecipeId(db *sql.DB, recipeId string) ([]Tag, error) {
	rows, err := db.Query(`
		SELECT id, recipe_id, tag
		FROM tags
		WHERE recipe_id = $1
	`, recipeId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag

	for rows.Next() {
		var tag Tag

		err := rows.Scan(
			&tag.ID,
			&tag.RecipeId,
			&tag.Tag,
		)

		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func addTagRoutes(rg *gin.RouterGroup) {
	tags := rg.Group("/tags")

	// Get all tags
	tags.GET("/", func(c *gin.Context) {

		db := c.MustGet("db").(*sql.DB)

		tags, err := getAllTags(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, tags)
	})
}
