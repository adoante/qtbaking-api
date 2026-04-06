package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func addVodRoutes(rg *gin.RouterGroup, db *sql.DB) {
	vods := rg.Group("/vods")

	// Get all vods
	// TODO: add pagination
	vods.GET("/", func(c *gin.Context) {
		rows, err := db.Query(
			`SELECT id, slug, title, video_url, created_at
			FROM vods
		`)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		// Return JSON response
		var vods []Vod

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
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			vods = append(vods, vod)
		}

		c.JSON(http.StatusOK, vods)
	})

	// get vod by slug
	vods.GET("/:slug", func(c *gin.Context) {
		id := c.Param("slug")

		rows, err := db.Query(
			`SELECT id, slug, title, video_url, created_at
			FROM vods
			WHERE slug = $1
		`, id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		// Return JSON response
		var vods []Vod

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
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			vods = append(vods, vod)
		}

		c.JSON(http.StatusOK, vods)
	})
}
