package dto

type AnalysisResponse struct {
	Latitude     float64                `json:"latitude" binding:"required"`
	Longitude    float64                `json:"longitude" binding:"required"`
	LocationName string                 `json:"location_name"`
	SoilData     map[string]interface{} `json:"soil_data"`
	WeatherData  map[string]interface{} `json:"weather_data"`
	Predictions  map[string]interface{} `json:"predictions"`
	AiResponse   string                 `json:"ai_response"`
}

type AnalysisRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

type MLServiceRequest struct {
	RawSoil    interface{} `json:"raw_soil"`
	RawWeather interface{} `json:"raw_weather"`
}
