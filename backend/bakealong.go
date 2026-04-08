package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func getBakealongByYtId(db *sql.DB, ytId string) (Bakealong, error) {
	var bakealong Bakealong

	// Get Vod Data
	vod, err := getVodBySlug(db, ytId)
	if err != nil {
		return bakealong, err
	}

	// Get Recipe Data
	var fullRecipes []RecipeResponse
	recipes, err := getRecipesByVodId(db, strconv.Itoa(vod.ID))

	if err != nil {
		return bakealong, err
	}

	// Get Components
	var fullComponents []ComponentResponse

	for _, recipe := range recipes {
		components, err := getComponentsByRecipeId(db, strconv.Itoa(recipe.ID))
		if err != nil {
			return bakealong, err
		}

		// Get Ingredients
		for _, component := range components {
			ingredients, err := getIngredientsByComponentId(db, strconv.Itoa(component.ID))
			if err != nil {
				return bakealong, err
			}

			fullComponents = append(fullComponents, ComponentResponse{
				Component:   component,
				Ingredients: ingredients,
			})
		}

		// Get Tools
		tools, err := getToolsByRecipeId(db, strconv.Itoa(recipe.ID))
		if err != nil {
			return bakealong, err
		}

		// Get Notes
		notes, err := getNotesByRecipeId(db, strconv.Itoa(recipe.ID))
		if err != nil {
			return bakealong, err
		}

		// Get Tags
		tags, err := getTagsByRecipeId(db, strconv.Itoa(recipe.ID))
		if err != nil {
			return bakealong, err
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

	bakealong.ID = vod.ID
	bakealong.Slug = vod.Slug
	bakealong.VodTitle = vod.Title
	bakealong.VideoURL = vod.VideoURL
	bakealong.CreatedAt = vod.CreatedAt

	bakealong.Recipes = fullRecipes

	return bakealong, nil
}

func addBakealongRoutes(rg *gin.RouterGroup) {
	bakealongs := rg.Group("/bakealong")

	bakealongs.GET("/:ytid", func(c *gin.Context) {
		db := c.MustGet("db").(*sql.DB)
		ytId := c.Param("ytid")

		bakealong, err := getBakealongByYtId(db, ytId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "bakealong not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, bakealong)
	})
}
