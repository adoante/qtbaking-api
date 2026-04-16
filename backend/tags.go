package main

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
	tags.GET("", func(c *gin.Context) {

		db := c.MustGet("db").(*sql.DB)

		filter := c.Query("filter")
		match := c.DefaultQuery("match", "exact")

		if match != "partial" && match != "exact" {
			match = "exact"
		}

		tags, err := getAllTags(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		var result []Tag
		if filter != "" && match == "exact" {
			for _, tag := range tags {
				if tag.Tag == filter {
					result = append(result, tag)
				}
			}

			c.JSON(http.StatusOK, result)
			return
		}

		if filter != "" && match == "partial" {
			for _, tag := range tags {
				if strings.Contains(tag.Tag, filter) {
					result = append(result, tag)
				}
			}

			c.JSON(http.StatusOK, result)
			return
		}

		c.JSON(http.StatusOK, tags)
	})
}
