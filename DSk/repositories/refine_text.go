package repositories

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

type CloudmersiveRephraseRequest struct {
	TextToRephrase string `json:"TextToRephrase"`
}

type CloudmersiveRephraseResponse struct {
	Successful       bool `json:"Successful"`
	RephrasedResults []struct {
		RephrasedSentence string `json:"RephrasedSentence"`
	} `json:"RephrasedResults"`
}

// RefineText sends the provided text to Cloudmersive's rephrase endpoint to refine it.
func RefineText(text string) (string, error) {
	apiKey := os.Getenv("CLOUDMERSIVE_API_KEY")
	if apiKey == "" {
		return "", errors.New("CLOUDMERSIVE_API_KEY not set in environment")
	}

	payload := CloudmersiveRephraseRequest{
		TextToRephrase: text,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.cloudmersive.com", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Apikey", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", errors.New(string(bodyBytes))
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response CloudmersiveRephraseResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return "", err
	}

	if response.Successful && len(response.RephrasedResults) > 0 {
		return response.RephrasedResults[0].RephrasedSentence, nil
	}

	return "", errors.New("failed to refine text")
}
