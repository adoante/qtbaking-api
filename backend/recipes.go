package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func getAllRecipes(db *sql.DB) ([]Recipe, error) {
	var recipes []Recipe

	rows, err := db.Query(
		`SELECT id, vod_id, title, thumbnail, temp_fahrenheit, temp_celsius
		FROM recipes
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var recipe Recipe

		err := rows.Scan(
			&recipe.ID,
			&recipe.VodId,
			&recipe.Title,
			&recipe.Thumbnail,
			&recipe.TempFahrenheit,
			&recipe.TempCelsius,
		)

		if err != nil {
			return nil, err
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func getAllFullRecipes(db *sql.DB) ([]FullRecipe, error) {
	// Get Recipe Data
	var fullRecipes []FullRecipe
	recipes, err := getAllRecipes(db)

	if err != nil {
		return fullRecipes, err
	}

	// Get Components
	for _, recipe := range recipes {
		var fullComponents []FullComponent

		components, err := getComponentsByRecipeId(db, strconv.Itoa(recipe.ID))
		if err != nil {
			return fullRecipes, err
		}

		// Get Ingredients
		for _, component := range components {
			ingredients, err := getIngredientsByComponentId(db, strconv.Itoa(component.ID))
			if err != nil {
				return fullRecipes, err
			}

			fullComponents = append(fullComponents, FullComponent{
				Component:   component,
				Ingredients: ingredients,
			})
		}

		// Get Tools
		tools, err := getToolsByRecipeId(db, strconv.Itoa(recipe.ID))
		if err != nil {
			return fullRecipes, err
		}

		// Get Notes
		notes, err := getNotesByRecipeId(db, strconv.Itoa(recipe.ID))
		if err != nil {
			return fullRecipes, err
		}

		// Get Tags
		tags, err := getTagsByRecipeId(db, strconv.Itoa(recipe.ID))
		if err != nil {
			return fullRecipes, err
		}

		fullRecipes = append(fullRecipes, FullRecipe{
			ID:             recipe.ID,
			Title:          recipe.Title,
			Thumbnail:      recipe.Thumbnail,
			TempFahrenheit: recipe.TempFahrenheit,
			TempCelsius:    recipe.TempCelsius,
			Components:     fullComponents,
			Tools:          tools,
			Notes:          notes,
			Tags:           tags,
		})
	}

	return fullRecipes, nil
}

func getFullRecipeId(db *sql.DB, id string) (FullRecipe, error) {
	// Get Recipe Data
	var fullRecipe FullRecipe
	recipe, err := getRecipeById(db, id)

	if err != nil {
		return fullRecipe, err
	}

	// Get Components
	var fullComponents []FullComponent

	components, err := getComponentsByRecipeId(db, strconv.Itoa(recipe.ID))
	if err != nil {
		return fullRecipe, err
	}

	// Get Ingredients
	for _, component := range components {
		ingredients, err := getIngredientsByComponentId(db, strconv.Itoa(component.ID))
		if err != nil {
			return fullRecipe, err
		}

		fullComponents = append(fullComponents, FullComponent{
			Component:   component,
			Ingredients: ingredients,
		})
	}

	// Get Tools
	tools, err := getToolsByRecipeId(db, strconv.Itoa(recipe.ID))
	if err != nil {
		return fullRecipe, err
	}

	// Get Notes
	notes, err := getNotesByRecipeId(db, strconv.Itoa(recipe.ID))
	if err != nil {
		return fullRecipe, err
	}

	// Get Tags
	tags, err := getTagsByRecipeId(db, strconv.Itoa(recipe.ID))
	if err != nil {
		return fullRecipe, err
	}

	fullRecipe = FullRecipe{
		ID:             recipe.ID,
		Title:          recipe.Title,
		Thumbnail:      recipe.Thumbnail,
		TempFahrenheit: recipe.TempFahrenheit,
		TempCelsius:    recipe.TempCelsius,
		Components:     fullComponents,
		Tools:          tools,
		Notes:          notes,
		Tags:           tags,
	}

	return fullRecipe, nil
}

func getRecipeById(db *sql.DB, id string) (Recipe, error) {
	var recipe Recipe

	err := db.QueryRow(`
		SELECT id, vod_id, title, thumbnail, temp_fahrenheit, temp_celsius
		FROM recipes
		WHERE id = $1
	`, id).Scan(
		&recipe.ID,
		&recipe.VodId,
		&recipe.Title,
		&recipe.Thumbnail,
		&recipe.TempFahrenheit,
		&recipe.TempCelsius,
	)

	if err != nil {
		return recipe, err
	}

	return recipe, nil
}

func getRecipesByVodId(db *sql.DB, vodId string) ([]Recipe, error) {
	var recipes []Recipe

	rows, err := db.Query(
		`SELECT id, vod_id, title, thumbnail, temp_fahrenheit, temp_celsius
		FROM recipes
		WHERE vod_id = $1
	`, vodId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var recipe Recipe

		err := rows.Scan(
			&recipe.ID,
			&recipe.VodId,
			&recipe.Title,
			&recipe.Thumbnail,
			&recipe.TempFahrenheit,
			&recipe.TempCelsius,
		)

		if err != nil {
			return nil, err
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil

}

func addRecipeRoutes(rg *gin.RouterGroup) {
	recipes := rg.Group("/recipes")

	recipes.GET("", func(c *gin.Context) {
		db := c.MustGet("db").(*sql.DB)

		filterTitle := c.Query("title")
		filterTag := c.Query("tag")
		match := c.DefaultQuery("match", "partial")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		if match != "partial" && match != "exact" {
			match = "exact"
		}

		recipes, err := getAllFullRecipes(db)

		end := min(offset+limit, len(recipes))
		start := min(offset, len(recipes))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		var result []FullRecipe

		if filterTag != "" && match == "exact" {
			for _, recipe := range recipes {
				for _, tag := range recipe.Tags {
					if tag.Tag == filterTag {
						result = append(result, recipe)
					}
				}
			}

			end := min(offset+limit, len(result))
			start := min(offset, len(result))

			result = result[start:end]
			c.JSON(http.StatusOK, result)
			return
		}

		if filterTag != "" && match == "partial" {
			for _, recipe := range recipes {
				for _, tag := range recipe.Tags {
					if strings.Contains(tag.Tag, filterTag) {
						result = append(result, recipe)
					}
				}
			}

			end := min(offset+limit, len(result))
			start := min(offset, len(result))

			result = result[start:end]
			c.JSON(http.StatusOK, result)
			return
		}

		if filterTitle != "" && match == "exact" {
			for _, recipe := range recipes {
				if recipe.Title == filterTitle {
					result = append(result, recipe)
				}
			}

			end := min(offset+limit, len(result))
			start := min(offset, len(result))

			result = result[start:end]
			c.JSON(http.StatusOK, result)
			return
		}

		if filterTitle != "" && match == "partial" {
			for _, recipe := range recipes {
				if strings.Contains(recipe.Title, filterTitle) {
					result = append(result, recipe)
				}
			}

			end := min(offset+limit, len(result))
			start := min(offset, len(result))

			result = result[start:end]
			c.JSON(http.StatusOK, result)
			return
		}

		c.JSON(http.StatusOK, recipes[start:end])
	})

	// Get recipe by id
	recipes.GET("/:id", func(c *gin.Context) {
		db := c.MustGet("db").(*sql.DB)

		id := c.Param("id")

		recipe, err := getFullRecipeId(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, recipe)
	})

}
