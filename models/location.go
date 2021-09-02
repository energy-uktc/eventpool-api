package models

type Location struct {
	Latitude    float32 `json:"latitude" binding:"required"`
	Longitude   float32 `json:"longitude" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
}
