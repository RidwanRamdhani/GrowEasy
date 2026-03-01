package integration

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

type GeminiClient struct {
	Client *genai.Client
	Model  string
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
		model = "gemini-3-flash-preview"
	}

	return &GeminiClient{
		Client: client,
		Model:  model,
	}, nil
}

// GenerateSummary creates a summary from ML prediction results using Gemini API
func (c *GeminiClient) GenerateSummary(
	predictionClass string,
	probability float64,
	top3 map[string]float64,
	soilData map[string]interface{},
	weatherData map[string]interface{},
) (string, error) {

	prompt := fmt.Sprintf(`Based on the following ML prediction results, provide a brief summary and agricultural recommendations with only these data:

PREDICTION RESULTS:
- Main Crop: %s
- Confidence: %.2f%%
- Top 3 Predictions: %v

SOIL CONDITIONS: %v

WEATHER CONDITIONS: %v

Please provide:
1. Brief analysis of the prediction
2. Agricultural recommendations
3. Risk assessment (if any)

Response in English, max 200 words.`,
		predictionClass,
		probability*100,
		top3,
		soilData,
		weatherData,
	)

	temp := float32(0.7)
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
