package models

type WeatherResponse struct {
	City        string `json:"city"`
	Temperature int    `json:"temperature"`
}
