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
	Weather interface{} `json:"weather"`
	Soil    interface{} `json:"soil"`
}

type Top3Prediction map[string]float64

type MLServiceResponse struct {
	PredictionClass string         `json:"prediction_class"`
	Probability     float64        `json:"probability"`
	Top3            Top3Prediction `json:"top3"`
}
