package main

import (
	"time"
)

type Vod struct {
	ID        int       `json:"id"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	VideoURL  string    `json:"video_url"`
	CreatedAt time.Time `json:"created_at"`
}

type Recipe struct {
	ID             int     `json:"id"`
	VodId          *int64  `json:"vod_id,omitempty"`
	Thumbnail      *string `json:"thumbnail,omitempty"`
	Title          string  `json:"title"`
	TempFahrenheit *int64  `json:"temp_fahrenheit,omitempty"`
	TempCelsius    *int64  `json:"temp_celsius,omitempty"`
}

type Component struct {
	ID       int    `json:"id"`
	RecipeId int    `json:"recipe_id"`
	Name     string `json:"name"`
}

type Ingredient struct {
	ID             int      `json:"id"`
	ComponentId    int      `json:"component_id"`
	Name           string   `json:"name"`
	Quantity       float64  `json:"quantity"`
	Unit           string   `json:"unit"`
	MetricQuantity *float64 `json:"metric_quantity,omitempty"`
	MetricUnit     *string  `json:"metric_unit,omitempty"`
	Optional       bool     `json:"optional"`
	Notes          *string  `json:"notes,omitempty"`
}

type Tool struct {
	ID       int    `json:"id"`
	RecipeId int    `json:"recipe_id"`
	Name     string `json:"name"`
	Optional bool   `json:"optional"`
}

type Note struct {
	ID       int    `json:"id"`
	RecipeId int    `json:"recipe_id"`
	Note     string `json:"note"`
}

type Tag struct {
	ID       int    `json:"id"`
	RecipeId int    `json:"recipe_id"`
	Tag      string `json:"tag"`
}

type FullComponent struct {
	Component
	Ingredients []Ingredient `json:"ingredients"`
}

type FullRecipe struct {
	ID             int             `json:"id"`
	Title          string          `json:"title"`
	Thumbnail      *string         `json:"thumbnail,omitempty"`
	TempFahrenheit *int64          `json:"temp_fahrenheit,omitempty"`
	TempCelsius    *int64          `json:"temp_celsius,omitempty"`
	Components     []FullComponent `json:"components"`
	Tools          []Tool          `json:"tools"`
	Notes          []Note          `json:"notes"`
	Tags           []Tag           `json:"tags"`
}

type Bakealong struct {
	ID        int          `json:"id"`
	Slug      string       `json:"slug"`
	VodTitle  string       `json:"vod_title"`
	VideoURL  string       `json:"video_url"`
	CreatedAt time.Time    `json:"created_at"`
	Recipes   []FullRecipe `json:"recipes"`
}
