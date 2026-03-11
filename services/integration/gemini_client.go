package integration

import (
	"context"
	"fmt"
	"os"
	"sync"

	"google.golang.org/genai"
)

type GeminiClient struct {
	Client   *genai.Client
	Model    string
	sessions map[string][]string // History per user
	mu       sync.RWMutex
}

func NewGeminiClient() (*GeminiClient, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-3.1-flash-lite-preview"
	}

	return &GeminiClient{
		Client:   client,
		Model:    model,
		sessions: make(map[string][]string),
	}, nil
}

// aggregateWeatherData computes summary stats from raw weather API response to avoid sending large raw arrays to Gemini.
func aggregateWeatherData(weatherData map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	daily, ok := weatherData["daily"].(map[string]interface{})
	if !ok {
		return weatherData
	}

	// Average temperature
	if temps, ok := daily["temperature_2m_mean"].([]interface{}); ok && len(temps) > 0 {
		sum := 0.0
		min := 1000.0
		max := -1000.0
		for _, v := range temps {
			if t, ok := v.(float64); ok {
				sum += t
				if t < min {
					min = t
				}
				if t > max {
					max = t
				}
			}
		}
		avg := sum / float64(len(temps))
		result["avg_temperature_c"] = fmt.Sprintf("%.1f", avg)
		result["min_temperature_c"] = fmt.Sprintf("%.1f", min)
		result["max_temperature_c"] = fmt.Sprintf("%.1f", max)
	}

	// Total & average daily precipitation
	if precips, ok := daily["precipitation_sum"].([]interface{}); ok && len(precips) > 0 {
		total := 0.0
		rainyDays := 0
		for _, v := range precips {
			if p, ok := v.(float64); ok {
				total += p
				if p > 1.0 {
					rainyDays++
				}
			}
		}
		result["total_precipitation_mm"] = fmt.Sprintf("%.1f", total)
		result["avg_daily_precipitation_mm"] = fmt.Sprintf("%.1f", total/float64(len(precips)))
		result["rainy_days_last_30d"] = rainyDays
	}

	// Average humidity from hourly data
	if hourly, ok := weatherData["hourly"].(map[string]interface{}); ok {
		if humidities, ok := hourly["relative_humidity_2m"].([]interface{}); ok && len(humidities) > 0 {
			sum := 0.0
			count := 0
			for _, v := range humidities {
				if h, ok := v.(float64); ok {
					sum += h
					count++
				}
			}
			if count > 0 {
				result["avg_humidity_pct"] = fmt.Sprintf("%.1f", sum/float64(count))
			}
		}
	}

	if tz, ok := weatherData["timezone"].(string); ok {
		result["timezone"] = tz
	}

	return result
}

// aggregateSoilData extracts and normalizes soil layer values,
func aggregateSoilData(soilData map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	properties, ok := soilData["properties"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{"note": "no soil data available"}
	}

	layers, ok := properties["layers"].([]interface{})
	if !ok {
		return map[string]interface{}{"note": "no soil layers available"}
	}

	for _, l := range layers {
		layer, ok := l.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := layer["name"].(string)
		unitMeasure, _ := layer["unit_measure"].(map[string]interface{})
		targetUnit, _ := unitMeasure["target_units"].(string)
		dFactor := 1.0
		if df, ok := unitMeasure["d_factor"].(float64); ok && df != 0 {
			dFactor = df
		}

		depths, ok := layer["depths"].([]interface{})
		if !ok {
			continue
		}

		layerValues := map[string]interface{}{}
		hasData := false

		for _, d := range depths {
			depth, ok := d.(map[string]interface{})
			if !ok {
				continue
			}

			label, _ := depth["label"].(string)
			values, _ := depth["values"].(map[string]interface{})
			mean := values["mean"]

			if mean == nil {
				layerValues[label] = "no data"
				continue
			}

			hasData = true
			if meanVal, ok := mean.(float64); ok {
				layerValues[label] = fmt.Sprintf("%.2f %s", meanVal/dFactor, targetUnit)
			}
		}

		if !hasData {
			result[name] = fmt.Sprintf("no data available (unit: %s)", targetUnit)
		} else {
			result[name] = layerValues
		}
	}

	return result
}

// GenerateSummary creates a summary from ML prediction results using Gemini API
func (c *GeminiClient) GenerateSummary(
	predictionClass string,
	probability float64,
	top3 map[string]float64,
	soilData map[string]interface{},
	weatherData map[string]interface{},
) (string, error) {
	weather := aggregateWeatherData(weatherData)
	soil := aggregateSoilData(soilData)

	// Format top3 as readable percentages
	top3Formatted := ""
	for crop, prob := range top3 {
		top3Formatted += fmt.Sprintf("\n   - %s: %.1f%%", crop, prob*100)
	}

	// Format soil as readable summary
	soilSummary := ""
	for k, v := range soil {
		soilSummary += fmt.Sprintf("\n   - %s: %v", k, v)
	}

	prompt := fmt.Sprintf(`You are an expert agronomist assistant. Based on the ML crop prediction and field conditions below, write a practical farming advisory for the farmer.

## ML PREDICTION
- Recommended Crop: %s (confidence: %.1f%%)
- Other Candidates:%s

## WEATHER CONDITIONS (last 30 days)
- Avg Temperature: %s°C (min %s°C / max %s°C)
- Total Rainfall: %s mm over 30 days (%s mm/day avg, %v rainy days)
- Avg Humidity: %s%%

## SOIL CONDITIONS%s

---

Write a farming advisory in English (200–250 words) covering ONLY these three sections:

**1. Fertilizer & Nutrient Recommendations**
Based on the soil nutrient profile above (nitrogen, pH), recommend appropriate fertilizer type, dosage guidance, and application timing. Give confident, direct recommendations — do NOT mention missing data, unavailable readings, or suggest the farmer conduct soil tests. Just provide the best general recommendation for the predicted crop and climate.

**2. Weather Risk Assessment**
Identify specific risks based on the 30-day weather data (e.g. waterlogging risk from high rainfall, heat stress, disease pressure from high humidity). Be specific with thresholds where relevant.

**3. Optimal Planting Schedule**
Suggest the best time window to plant %s given current conditions, including any preparation steps before planting. If conditions are not yet ideal, state what the farmer should wait for.

Keep the tone practical and direct — this is for a farmer in the field, not an academic.`,
		predictionClass,
		probability*100,
		top3Formatted,
		weather["avg_temperature_c"], weather["min_temperature_c"], weather["max_temperature_c"],
		weather["total_precipitation_mm"], weather["avg_daily_precipitation_mm"], weather["rainy_days_last_30d"],
		weather["avg_humidity_pct"],
		soilSummary,
		predictionClass,
	)

	temp := float32(0.5)
	result, err := c.Client.Models.GenerateContent(
		context.Background(),
		c.Model,
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			Temperature: &temp,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(result.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}

	candidate := result.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("no content in candidate")
	}

	summary := candidate.Content.Parts[0].Text

	return summary, nil
}

// ChatWithSession handles chat with in-memory history per user
func (c *GeminiClient) ChatWithSession(userID, analysisContext, message string) (string, error) {
	c.mu.Lock()
	history, exists := c.sessions[userID]
	if !exists {
		// Initialize with system prompt
		systemPrompt := fmt.Sprintf(`You are an expert agronomist assistant. Use the following latest analysis context to inform your responses:

%s

Engage in a conversational manner with the farmer. Answer questions based on this context and general farming knowledge. Keep responses concise and brief`, analysisContext)
		history = []string{systemPrompt}
		c.sessions[userID] = history
	}
	c.mu.Unlock()

	// Add user message to history
	c.mu.Lock()
	c.sessions[userID] = append(c.sessions[userID], "User: "+message)
	history = c.sessions[userID]
	c.mu.Unlock()

	// Build full prompt
	fullPrompt := ""
	for _, h := range history {
		fullPrompt += h + "\n"
	}

	temp := float32(0.)
	result, err := c.Client.Models.GenerateContent(
		context.Background(),
		c.Model,
		genai.Text(fullPrompt),
		&genai.GenerateContentConfig{
			Temperature: &temp,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate chat response: %w", err)
	}

	if len(result.Candidates) == 0 || result.Candidates[0].Content == nil || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	response := result.Candidates[0].Content.Parts[0].Text

	// Add AI response to history
	c.mu.Lock()
	c.sessions[userID] = append(c.sessions[userID], "AI: "+response)
	c.mu.Unlock()

	return response, nil
}

// ClearSession clears the session for a user
func (c *GeminiClient) ClearSession(userID string) {
	c.mu.Lock()
	delete(c.sessions, userID)
	c.mu.Unlock()
}
