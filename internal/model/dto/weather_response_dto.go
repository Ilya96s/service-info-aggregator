package dto

type WeatherResponse struct {
	City string `json:"city"`
	Temp int    `json:"temp"`
}
