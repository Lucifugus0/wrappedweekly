package aiprovider

import "wrappedweekly/backend/internal/domain"

// NewProvider selects an AIProvider implementation based on the AI_PROVIDER
// env var. Only "mock" is implemented out of the box (deterministic, no API
// key needed). To wire a real LLM: implement domain.AIProvider (e.g. a new
// OpenAIProvider/AnthropicProvider struct calling the actual API) and add a
// case for it here, then set AI_PROVIDER=<name> and the relevant API key env var.
func NewProvider(name string) domain.AIProvider {
	switch name {
	case "mock":
		return NewMockProvider()
	default:
		return NewMockProvider()
	}
}
