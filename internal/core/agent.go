package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/myuser/go-agent/internal/channels"
	"github.com/myuser/go-agent/internal/config"
	"github.com/myuser/go-agent/internal/llm"
	"github.com/myuser/go-agent/internal/tools"
	"github.com/sashabaranov/go-openai"
)

type Agent struct {
	cfg          *config.Config
	llm          llm.Provider
	channels     *channels.Manager
	tools        map[string]tools.Tool
	history      []openai.ChatCompletionMessage
	systemPrompt string
}

func NewAgent(cfg *config.Config, llm llm.Provider, cm *channels.Manager) *Agent {
	return &Agent{
		cfg:      cfg,
		llm:      llm,
		channels: cm,
		tools:    make(map[string]tools.Tool),
		history:  make([]openai.ChatCompletionMessage, 0),
	}
}

func (a *Agent) RegisterTool(t tools.Tool) {
	a.tools[t.Name()] = t
}

func (a *Agent) SetSystemPrompt(prompt string) {
	a.systemPrompt = prompt
	// Reset history with system prompt
	a.history = []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: prompt,
		},
	}
}

func (a *Agent) Run() {
	go a.channels.Start()

	log.Println("Agent started. Waiting for messages...")
	for msg := range a.channels.Messages() {
		go a.handleMessage(msg)
	}
}

func (a *Agent) handleMessage(msg channels.Message) {
	log.Printf("Received message from %s: %s", msg.Sender, msg.Content)

	// Add user message to history
	a.history = append(a.history, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: msg.Content,
	})

	// Prepare tools for LLM
	llmTools := []openai.Tool{}
	for _, t := range a.tools {
		schema := t.Schema()
		llmTools = append(llmTools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        t.Name(),
				Description: t.Description(),
				Parameters:  schema,
			},
		})
	}

	ctx := context.Background()

	// Loop to handle tool calls
	maxTurns := 5
	for i := 0; i < maxTurns; i++ {
		// Show thinking indicator
		a.channels.ShowThinking(msg.Channel)

		// Use ChatStream for streaming response
		stream, err := a.llm.ChatStream(ctx, a.history, llmTools)
		if err != nil {
			log.Printf("LLM error: %v", err)
			a.channels.SendToChannel(msg.Channel, "Error communicating with AI.")
			return
		}
		defer stream.Close()

		var contentBuilder string
		var toolCalls []openai.ToolCall
		isStreamable := a.channels.IsChannelStreamable(msg.Channel)
		// We don't need to send "[Agent]: " prefix here because the UI handles bubble creation.
		// Sending it causes "Agent: " to appear inside the bubble or multiple bubbles if loop repeats.
		// if isStreamable {
		// 	a.channels.SendTokenToChannel(msg.Channel, "[Agent]: ")
		// }

		for {
			response, err := stream.Recv()
			if err != nil {
				// End of stream or error
				break
			}

			if len(response.Choices) > 0 {
				delta := response.Choices[0].Delta
				if delta.Content != "" {
					contentBuilder += delta.Content
					if isStreamable {
						a.channels.SendTokenToChannel(msg.Channel, delta.Content)
					}
				}

				// Support for Reasoning Content (DeepSeek R1 style)
				// Note: This field is available in recent go-openai versions
				if rc := delta.ReasoningContent; rc != "" {
					if isStreamable {
						a.channels.SendReasoningToChannel(msg.Channel, rc)
					}
				}

				if len(delta.ToolCalls) > 0 {
					// Append delta tool calls
					// Note: Azure/OpenAI returns tool calls chunks. We need to assemble them.
					// This logic can be complex. For simplicity, we assume we need to accumulate them by index.
					for _, tcDelta := range delta.ToolCalls {
						if tcDelta.Index != nil {
							idx := *tcDelta.Index
							// Resize slice if needed
							if idx >= len(toolCalls) {
								newToolCalls := make([]openai.ToolCall, idx+1)
								copy(newToolCalls, toolCalls)
								toolCalls = newToolCalls
							}

							if tcDelta.ID != "" {
								toolCalls[idx].ID = tcDelta.ID
								toolCalls[idx].Type = tcDelta.Type
							}
							if tcDelta.Function.Name != "" {
								toolCalls[idx].Function.Name += tcDelta.Function.Name
							}
							if tcDelta.Function.Arguments != "" {
								toolCalls[idx].Function.Arguments += tcDelta.Function.Arguments
							}
						}
					}
				}
			}
		}

		if isStreamable {
			a.channels.SendTokenToChannel(msg.Channel, "\n")
		}

		// Construct the complete message
		msgResp := openai.ChatCompletionMessage{
			Role:      openai.ChatMessageRoleAssistant,
			Content:   contentBuilder,
			ToolCalls: toolCalls,
		}
		a.history = append(a.history, msgResp)

		// Check for tool calls
		if len(toolCalls) > 0 {
			for _, toolCall := range toolCalls {
				log.Printf("Tool call: %s %s", toolCall.Function.Name, toolCall.Function.Arguments)

				// Notify UI about tool execution
				if isStreamable {
					a.channels.SendToolCallToChannel(msg.Channel, toolCall.Function.Name, toolCall.Function.Arguments)
				}

				toolName := toolCall.Function.Name
				tool, exists := a.tools[toolName]
				if !exists {
					log.Printf("Tool not found: %s", toolName)
					a.history = append(a.history, openai.ChatCompletionMessage{
						Role:       openai.ChatMessageRoleTool,
						Content:    fmt.Sprintf("Error: Tool %s not found", toolName),
						ToolCallID: toolCall.ID,
					})
					continue
				}

				result, err := tool.Execute(json.RawMessage(toolCall.Function.Arguments))
				if err != nil {
					result = fmt.Sprintf("Error: %v", err)
					log.Printf("Tool error: %v", err)
				} else {
					log.Printf("Tool output: %s", result)
				}

				a.history = append(a.history, openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    result,
					ToolCallID: toolCall.ID,
				})
			}
			// Continue loop to send tool outputs back to LLM
		} else {
			// Final response
			if !isStreamable {
				a.channels.SendToChannel(msg.Channel, contentBuilder)
			}
			return
		}
	}
}
