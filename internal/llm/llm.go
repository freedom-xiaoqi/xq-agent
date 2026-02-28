package llm

import (
	"context"

	"xq-agent/internal/config"
	"github.com/sashabaranov/go-openai"
)

type Provider interface {
	Chat(ctx context.Context, messages []openai.ChatCompletionMessage, tools []openai.Tool) (openai.ChatCompletionResponse, error)
	ChatStream(ctx context.Context, messages []openai.ChatCompletionMessage, tools []openai.Tool) (*openai.ChatCompletionStream, error)
}

type OpenAIProvider struct {
	client *openai.Client
	model  string
}

func NewOpenAI(cfg config.LLMConfig) *OpenAIProvider {
	c := openai.DefaultConfig(cfg.APIKey)
	if cfg.BaseURL != "" {
		c.BaseURL = cfg.BaseURL
	}
	return &OpenAIProvider{
		client: openai.NewClientWithConfig(c),
		model:  cfg.Model,
	}
}

func (p *OpenAIProvider) Chat(ctx context.Context, messages []openai.ChatCompletionMessage, tools []openai.Tool) (openai.ChatCompletionResponse, error) {
	req := openai.ChatCompletionRequest{
		Model:    p.model,
		Messages: messages,
	}
	if len(tools) > 0 {
		req.Tools = tools
	}
	return p.client.CreateChatCompletion(ctx, req)
}

func (p *OpenAIProvider) ChatStream(ctx context.Context, messages []openai.ChatCompletionMessage, tools []openai.Tool) (*openai.ChatCompletionStream, error) {
	req := openai.ChatCompletionRequest{
		Model:    p.model,
		Messages: messages,
		Stream:   true,
	}
	if len(tools) > 0 {
		req.Tools = tools
	}
	return p.client.CreateChatCompletionStream(ctx, req)
}
