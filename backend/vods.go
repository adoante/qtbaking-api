package main

import (
	"database/sql"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func getAllVods(db *sql.DB) ([]Vod, error) {
	var vods []Vod

	rows, err := db.Query(`
		SELECT id, slug, title, video_url, created_at
		FROM vods
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var vod Vod

		err := rows.Scan(
			&vod.ID,
			&vod.Slug,
			&vod.Title,
			&vod.VideoURL,
			&vod.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		vods = append(vods, vod)
	}

	return vods, nil
}

func getVodBySlug(db *sql.DB, slug string) (Vod, error) {
	var vod Vod

	err := db.QueryRow(`
		SELECT id, slug, title, video_url, created_at
		FROM vods
		WHERE slug = $1
	`, slug).Scan(
		&vod.ID,
		&vod.Slug,
		&vod.Title,
		&vod.VideoURL,
		&vod.CreatedAt,
	)

	if err != nil {
		return vod, err
	}

	return vod, nil
}

func addVodRoutes(rg *gin.RouterGroup) {
	vods := rg.Group("/vods")

	// Get all vods
	// /vods?sort=created_at&order=asc
	vods.GET("/", func(c *gin.Context) {
		db := c.MustGet("db").(*sql.DB)

		// https://gin-gonic.com/en/docs/routing/api-design/#filtering-and-sorting
		sortBy := c.DefaultQuery("sort", "created_at")
		order := c.DefaultQuery("order", "desc")

		// Validate the sort field against an allow-list to prevent injection.
		allowed := map[string]bool{"created_at": true}
		if !allowed[sortBy] {
			sortBy = "created_at"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		vods, err := getAllVods(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if order == "asc" && sortBy == "created_at" {
			sort.Slice(vods, func(a, b int) bool {
				return vods[a].CreatedAt.Before(vods[b].CreatedAt)
			})
		}

		if order == "desc" && sortBy == "created_at" {
			sort.Slice(vods, func(a, b int) bool {
				return vods[a].CreatedAt.After(vods[b].CreatedAt)
			})
		}

		c.JSON(http.StatusOK, vods)
	})

	// get vod by slug
	vods.GET("/:slug", func(c *gin.Context) {
		db := c.MustGet("db").(*sql.DB)
		slug := c.Param("slug")

		vod, err := getVodBySlug(db, slug)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "vod not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, vod)
	})
}
