package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
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
	vods.GET("/", func(c *gin.Context) {
		db := c.MustGet("db").(*sql.DB)

		vods, err := getAllVods(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
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
