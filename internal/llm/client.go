package llm

import (
	"context"
	"fmt"
	"time"
)

// Client defines the interface for an LLM service
type Client interface {
	Query(ctx context.Context, input string) (string, error)
}

// StubClient is a fake LLM implementation for testing
type StubClient struct{}

// NewStubClient returns a new StubClient
func NewStubClient() *StubClient {
	return &StubClient{}
}

// Query returns a canned response with simulated processing delay
func (c *StubClient) Query(ctx context.Context, input string) (string, error) {
	// Simulate some processing time
	select {
	case <-time.After(100 * time.Millisecond):
		// For example, if input contains "2 + 2", return "4"
		if input == "What is 2 + 2?" {
			return "4", nil
		}
		return fmt.Sprintf("stub response for: %s", input), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

