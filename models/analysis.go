package models

import (
	"time"

	"gorm.io/datatypes"
)

type Analysis struct {
	ID           string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID       *string        `gorm:"type:uuid;index" json:"user_id,omitempty"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Latitude     float64        `gorm:"not null" json:"latitude"`
	Longitude    float64        `gorm:"not null" json:"longitude"`
	LocationName string         `gorm:"type:text" json:"location_name,omitempty"`
	SoilData     datatypes.JSON `gorm:"type:jsonb" json:"soil_data,omitempty"`
	WeatherData  datatypes.JSON `gorm:"type:jsonb" json:"weather_data,omitempty"`
	Predictions  datatypes.JSON `gorm:"type:jsonb" json:"predictions,omitempty"`
	AiResponse   string         `gorm:"type:text" json:"ai_response,omitempty"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
}

func (Analysis) TableName() string {
	return "analyses"
}
