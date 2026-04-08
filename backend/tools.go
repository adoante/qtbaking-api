package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAllTools(db *sql.DB) ([]Tool, error) {
	rows, err := db.Query(`
		SELECT id, recipe_id, name, Optional
		FROM tools
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
			return nil, err
		}

		tools = append(tools, tool)
	}

	return tools, nil
}

func getToolsByRecipeId(db *sql.DB, recipeId string) ([]Tool, error) {
	rows, err := db.Query(`
		SELECT id, recipe_id, name, Optional
		FROM tools
		WHERE recipe_id = $1
	`, recipeId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
			return nil, err
		}

		tools = append(tools, tool)
	}

	return tools, nil
}

func addToolRoutes(rg *gin.RouterGroup) {
	tools := rg.Group("/tools")

	// Get all notes
	tools.GET("/", func(c *gin.Context) {

		db := c.MustGet("db").(*sql.DB)

		tools, err := getAllTools(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, tools)
	})
}
