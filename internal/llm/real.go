package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type RealClient struct {
	baseURL string
	client  *http.Client
}

func NewRealClient(baseURL string) *RealClient {
	return &RealClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	MaxTokens int          `json:"max_tokens"`
	Temperature float64	   `json:"temperature"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *RealClient) Query(ctx context.Context, input string) (string, error) {
	reqBody := chatRequest{
		Model: "tinyllama",
		Messages: []chatMessage{
			{Role: "system", Content: "You are a calculator. Always respond in a single line in the format <numeric expression> = <numeric answer>. Example Response: 2+2=4"},
			{Role: "user", Content: input},
		},
		MaxTokens: 50,
		Temperature: 0.1,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/v1/chat/completions",
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("llama server returned non-200 status")
	}

	var parsed chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}

	if len(parsed.Choices) == 0 {
		return "", errors.New("no choices returned from LLM")
	}

	return parsed.Choices[0].Message.Content, nil
}