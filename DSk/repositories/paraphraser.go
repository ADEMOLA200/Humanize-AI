package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ParaphraseResponse struct {
	Paraphrased string `json:"paraphrased"`
	Success     bool   `json:"success"`
}

func ParaphraseText(text string) (string, error) {
	requestBody, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://localhost:5001/paraphrase", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Uncomment the next line to see the raw response for debugging
	fmt.Println("Raw Response:", string(body))

	var result ParaphraseResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if !result.Success {
		return "", fmt.Errorf("failed to paraphrase")
	}

	return result.Paraphrased, nil
}
