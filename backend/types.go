package main

import (
	"database/sql"
	"time"
)

type Vod struct {
	ID        int
	Slug      string
	Title     string
	VideoURL  string
	CreatedAt time.Time
}

type Recipe struct {
	ID             int
	VodId          sql.NullInt64
	Thumbnail      sql.NullString
	Title          string
	TempFahrenheit sql.NullInt64
	TempCelsius    sql.NullInt64
}

type Component struct {
	ID       int
	RecipeId int
	Name     string
}

type Ingredient struct {
	ID             int
	ComponentId    int
	Name           string
	Quantity       float64
	Unit           string
	MetricQuantity sql.NullFloat64
	MetricUnit     sql.NullString
	Optional       bool
	Notes          string
}

type Tool struct {
	ID       int
	RecipeId int
	Name     string
	Optional bool
}

type Note struct {
	ID       int
	RecipeId int
	Note     string
}

type Tag struct {
	ID       int
	RecipeId int
	Tag      string
}
