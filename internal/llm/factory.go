package llm

import (
	"github.com/Octrafic/octrafic-cli/internal/llm/claude"
	"github.com/Octrafic/octrafic-cli/internal/llm/common"
	"github.com/Octrafic/octrafic-cli/internal/llm/openai"
	"fmt"
)

// CreateProvider creates a provider based on the config
func CreateProvider(config common.ProviderConfig) (common.Provider, error) {
	switch config.Provider {
	case "claude", "anthropic":
		return claude.NewClaudeProvider(config)
	case "openai", "openrouter":
		return openai.NewOpenAIProvider(config)
	default:
		return nil, fmt.Errorf("unsupported provider: %s (supported: claude, anthropic, openai, openrouter)", config.Provider)
	}
}
