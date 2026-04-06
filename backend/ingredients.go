package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func addIngredientsRoutes(rg *gin.RouterGroup, db *sql.DB) {
	ingredients := rg.Group("/ingredients")

	// Get all ingredients
	ingredients.GET("/", func(c *gin.Context) {
		rows, err := db.Query(
			`SELECT id, component_id, name, quantity, unit, metric_quantity, metric_unit, optional, notes
			FROM ingredients
		`)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		// Return JSON response
		var ingredients []Ingredient

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
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			ingredients = append(ingredients, ingredient)
		}

		c.JSON(http.StatusOK, ingredients)
	})
}
