package converter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ConversionRequest holds the input data from the user.
type ConversionRequest struct {
	Amount       float64
	FromCurrency string
	ToCurrency   string
}

// APIResponse maps the JSON structure returned by the exchange rate service.
type APIResponse struct {
	Result         string  `json:"result"`
	ConversionRate float64 `json:"conversion_rate"`
}

// GetRate performs an HTTP GET request to fetch the current exchange rate.
func GetRate(apiKey, from, to string) (float64, error) {
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/pair/%s/%s", apiKey, from, to)

	// Set a timeout to ensure the server doesn't hang on slow network responses
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var result APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.ConversionRate, nil
}

// Convert processes the currency calculation.
func Convert(req ConversionRequest, rate float64) float64 {
	return req.Amount * rate
}
