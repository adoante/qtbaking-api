package main

import (
	"database/sql"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func getBakealongByYtId(db *sql.DB, ytId string) (Bakealong, error) {
	var bakealong Bakealong

	// Get Vod Data
	vod, err := getVodBySlug(db, ytId)
	if err != nil {
		return bakealong, err
	}

	// Get Recipe Data
	var fullRecipes []FullRecipe
	recipes, err := getRecipesByVodId(db, strconv.Itoa(vod.ID))

	if err != nil {
		return bakealong, err
	}

	// Get Components
	var fullComponents []FullComponent

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

			fullComponents = append(fullComponents, FullComponent{
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

	bakealong.ID = vod.ID
	bakealong.Slug = vod.Slug
	bakealong.VodTitle = vod.Title
	bakealong.VideoURL = vod.VideoURL
	bakealong.CreatedAt = vod.CreatedAt

	bakealong.Recipes = fullRecipes

	return bakealong, nil
}

func getAllBakealongs(db *sql.DB) ([]Bakealong, error) {
	var bakealongs []Bakealong

	// Get all vod ids
	vods, err := getAllVods(db)
	if err != nil {
		return bakealongs, nil
	}

	for _, vod := range vods {
		bakealong, err := getBakealongByYtId(db, vod.Slug)
		if err != nil {
			return bakealongs, nil
		}
		bakealongs = append(bakealongs, bakealong)
	}

	return bakealongs, nil
}

func addBakealongRoutes(rg *gin.RouterGroup) {
	bakealongs := rg.Group("/bakealongs")

	bakealongs.GET("", func(c *gin.Context) {
		db := c.MustGet("db").(*sql.DB)

		sortBy := c.DefaultQuery("sort", "created_at")
		order := c.DefaultQuery("order", "desc")
		filterTags := c.QueryArray("tag")
		filterTitle := c.Query("title")
		match := c.DefaultQuery("match", "partial")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		if match != "partial" && match != "exact" {
			match = "exact"
		}

		if limit < 0 {
			limit = 10
		}
		if offset < 0 {
			offset = 0
		}

		allowed := map[string]bool{"created_at": true}
		if !allowed[sortBy] {
			sortBy = "created_at"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		bakealongs, err := getAllBakealongs(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if order == "asc" && sortBy == "created_at" {
			sort.Slice(bakealongs, func(a, b int) bool {
				return bakealongs[a].CreatedAt.Before(bakealongs[b].CreatedAt)
			})
		}

		if order == "desc" && sortBy == "created_at" {
			sort.Slice(bakealongs, func(a, b int) bool {
				return bakealongs[a].CreatedAt.After(bakealongs[b].CreatedAt)
			})
		}

		paginate := func(items []Bakealong) []Bakealong {
			if offset >= len(items) {
				return []Bakealong{}
			}

			end := offset + limit
			if end > len(items) {
				end = len(items)
			}

			return items[offset:end]
		}

		var result []Bakealong

		if match == "exact" {
			for _, bakealong := range bakealongs {
				matchesTitle := true
				matchesTags := true

				for _, recipe := range bakealong.Recipes {
					if len(filterTags) != 0 {
						var tags []string
						for _, tag := range recipe.Tags {
							tags = append(tags, tag.Tag)
						}
						matchesTags = containsAll(tags, filterTags)
					}
				}

				if filterTitle != "" {
					matchesTitle = bakealong.VodTitle == filterTitle
				}

				if matchesTags && matchesTitle {
					result = append(result, bakealong)
				}
			}

			c.JSON(http.StatusOK, paginate(result))
			return
		}

		if match == "partial" {
			for _, bakealong := range bakealongs {
				matchesTags := true
				matchesTitle := true

				for _, recipe := range bakealong.Recipes {
					if len(filterTags) != 0 {
						var tags []string
						for _, tag := range recipe.Tags {
							tags = append(tags, tag.Tag)
						}
						matchesTags = containsAny(tags, filterTags)
					}
				}

				if filterTitle != "" {
					matchesTitle = strings.Contains(bakealong.VodTitle, filterTitle)
				}

				if matchesTags && matchesTitle {
					result = append(result, bakealong)
				}
			}

			c.JSON(http.StatusOK, paginate(result))
			return
		}

		c.JSON(http.StatusOK, paginate(bakealongs))
	})

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
