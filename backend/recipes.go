package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
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

	// Get recipe by id
	recipes.GET("/", func(c *gin.Context) {
		db := c.MustGet("db").(*sql.DB)

		recipes, err := getAllRecipes(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, recipes)
	})

	// Get recipe by id
	recipes.GET("/:id", func(c *gin.Context) {
		db := c.MustGet("db").(*sql.DB)

		id := c.Param("id")

		recipe, err := getRecipeById(db, id)
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
