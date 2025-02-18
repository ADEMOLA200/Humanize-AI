package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ParaphraseResponse struct {
	Paraphrased string `json:"paraphrased"`
	Success     bool   `json:"success"`
}

func ParaphraseText(text string) (string, error) {
	// Load API URL from environment variable, default to localhost
	baseURL, exists := os.LookupEnv("PARAPHRASE_API_URL")
	if !exists {
		baseURL = "http://localhost:5001"
	}
	fmt.Printf("Using Paraphrase API: %s\n", baseURL)

	requestBody, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/paraphrase", baseURL), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	fmt.Println("Raw Response:", string(body))

	var result ParaphraseResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error parsing JSON response: %w", err)
	}

	if !result.Success {
		return "", fmt.Errorf("paraphrasing failed: %s", string(body))
	}

	return result.Paraphrased, nil
}
