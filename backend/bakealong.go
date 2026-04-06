package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func addBakealongRoutes(rg *gin.RouterGroup, db *sql.DB) {
	bakealongs := rg.Group("/bakealong")

	bakealongs.GET("/:ytid", func(c *gin.Context) {
		YtId := c.Param("ytid")

		// Get Vod data
		var vodRecipe VodRecipeResponse

		rows, err := db.Query(
			`SELECT id, slug, title, video_url, created_at
			FROM vods
			WHERE slug = $1
		`, YtId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		defer rows.Close()

		var vod Vod

		for rows.Next() {
			err := rows.Scan(
				&vod.ID,
				&vod.Slug,
				&vod.Title,
				&vod.VideoURL,
				&vod.CreatedAt,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		vodRecipe.ID = vod.ID
		vodRecipe.Slug = vod.Slug
		vodRecipe.VodTitle = vod.Title
		vodRecipe.VideoURL = vod.VideoURL
		vodRecipe.CreatedAt = vod.CreatedAt

		// Get Recipe data
		var fullRecipes []RecipeResponse

		rows, err = db.Query(
			`SELECT id, vod_id, thumbnail, title, temp_fahrenheit, temp_celsius
			FROM recipes
			WHERE vod_id = $1
		`, vod.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		var recipes []Recipe

		for rows.Next() {
			var recipe Recipe
			err := rows.Scan(
				&recipe.ID,
				&recipe.VodId,
				&recipe.Thumbnail,
				&recipe.Title,
				&recipe.TempFahrenheit,
				&recipe.TempCelsius,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			recipes = append(recipes, recipe)
		}

		// Get components
		var fullComponents []ComponentResponse

		for _, recipe := range recipes {
			var components []Component

			rows, err = db.Query(
				`SELECT id, recipe_id, name
				FROM components
				WHERE recipe_id = $1
				`, recipe.ID)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			defer rows.Close()

			for rows.Next() {
				var component Component
				err := rows.Scan(
					&component.ID,
					&component.RecipeId,
					&component.Name,
				)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				components = append(components, component)
			}

			// Get ingredients
			for _, component := range components {
				var ingredients []Ingredient

				rows, err = db.Query(
					`SELECT id, component_id, name, quantity, unit, metric_quantity, metric_unit, optional, notes
					FROM ingredients
					WHERE component_id = $1
					`, component.ID)

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
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
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					ingredients = append(ingredients, ingredient)
				}

				fullComponents = append(fullComponents, ComponentResponse{Component: component, Ingredients: ingredients})
			}

			// Get Tools
			rows, err := db.Query(
				`SELECT id, recipe_id, name, optional
				FROM tools
				WHERE recipe_id = $1
				`, recipe.ID)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
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
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				tools = append(tools, tool)
			}

			// Get Notes
			rows, err = db.Query(
				`SELECT id, recipe_id, note
				FROM notes
				WHERE recipe_id = $1
				`, recipe.ID)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
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
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				notes = append(notes, note)
			}

			// Get Tags
			rows, err = db.Query(
				`SELECT id, recipe_id, tag
				FROM tags
				WHERE recipe_id = $1
				`, recipe.ID)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
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
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				tags = append(tags, tag)
			}

			fullRecipes = append(fullRecipes, RecipeResponse{
				ID:             recipe.ID,
				RecipeTitle:    recipe.Title,
				Thumbnail:      recipe.Thumbnail,
				TempFahrenheit: recipe.TempFahrenheit,
				TempCelsius:    recipe.TempCelsius,
				Components:     fullComponents,
				Tools:          tools,
				Notes:          notes,
				Tags:           tags,
			})
		}

		vodRecipe.Recipes = fullRecipes

		c.JSON(http.StatusOK, vodRecipe)
	})
}
