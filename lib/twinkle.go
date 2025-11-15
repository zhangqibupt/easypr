package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TwinkleRequest struct {
	Prompt string `json:"prompt"`
}

type TwinkleResponse struct {
	Completion string `json:"completion"`
}

func SendTwinkle(ctx context.Context, msg string) (string, error) {
	// TODO need to set the url from config
	url := "xxx"

	requestBody, err := json.Marshal(TwinkleRequest{
		Prompt: msg,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status: %s, body: %s", resp.Status, string(body))
	}

	var twinkleResponse TwinkleResponse
	if err := json.Unmarshal(body, &twinkleResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return twinkleResponse.Completion, nil
}
