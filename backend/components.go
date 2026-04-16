package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAllComponents(db *sql.DB) ([]Component, error) {
	rows, err := db.Query(`
		SELECT id, recipe_id, name
		FROM components
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var components []Component

	for rows.Next() {
		var component Component

		err := rows.Scan(
			&component.ID,
			&component.RecipeId,
			&component.Name,
		)

		if err != nil {
			return nil, err
		}

		components = append(components, component)
	}

	return components, nil
}

func getComponentsByRecipeId(db *sql.DB, recipeId string) ([]Component, error) {
	rows, err := db.Query(`
		SELECT id, recipe_id, name
		FROM components
		WHERE recipe_id = $1
	`, recipeId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var components []Component

	for rows.Next() {
		var component Component

		err := rows.Scan(
			&component.ID,
			&component.RecipeId,
			&component.Name,
		)

		if err != nil {
			return nil, err
		}

		components = append(components, component)
	}

	return components, nil
}

func getComponentById(db *sql.DB, id string) (Component, error) {
	var component Component

	err := db.QueryRow(`
		SELECT id, recipe_id, name
		FROM components
		WHERE id = $1
	`, id).Scan(
		&component.ID,
		&component.RecipeId,
		&component.Name,
	)

	if err != nil {
		return component, err
	}

	return component, nil
}

func addComponentRoutes(rg *gin.RouterGroup) {
	components := rg.Group("/components")

	// Get all components
	components.GET("", func(c *gin.Context) {
		db := c.MustGet("db").(*sql.DB)

		components, err := getAllComponents(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, components)
	})
}
