package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAllIngredients(db *sql.DB) ([]Ingredient, error) {
	var ingredients []Ingredient

	rows, err := db.Query(`
		SELECT id, component_id, name, quantity, unit, metric_quantity, metric_unit, optional, notes
		FROM ingredients
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ingredient Ingredient

		err := rows.Scan(
			&ingredient.ID,
			&ingredient.ComponentId,
			&ingredient.Name,
			&ingredient.Quantity,
			&ingredient.Unit,
			&ingredient.MetricQuantity,
			&ingredient.MetricUnit,
			&ingredient.Optional,
			&ingredient.Notes,
		)

		if err != nil {
			return nil, err
		}

		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

func getIngredientsByComponentId(db *sql.DB, componentId string) ([]Ingredient, error) {
	var ingredients []Ingredient

	rows, err := db.Query(`
		SELECT id, component_id, name, quantity, unit, metric_quantity, metric_unit, optional, notes
		FROM ingredients
		WHERE component_id = $1
	`, componentId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ingredient Ingredient

		err := rows.Scan(
			&ingredient.ID,
			&ingredient.ComponentId,
			&ingredient.Name,
			&ingredient.Quantity,
			&ingredient.Unit,
			&ingredient.MetricQuantity,
			&ingredient.MetricUnit,
			&ingredient.Optional,
			&ingredient.Notes,
		)

		if err != nil {
			return nil, err
		}

		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

func addIngredientRoutes(rg *gin.RouterGroup) {
	ingredients := rg.Group("/ingredients")

	// Get all ingredients
	ingredients.GET("", func(c *gin.Context) {
		db := c.MustGet("db").(*sql.DB)

		ingredients, err := getAllIngredients(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, ingredients)
	})
}
