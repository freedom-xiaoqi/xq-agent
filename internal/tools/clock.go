package tools

import (
	"encoding/json"
	"time"
)

type ClockTool struct{}

func (t *ClockTool) Name() string { return "clock_current" }
func (t *ClockTool) Description() string {
	return "Get the current date and time."
}

func (t *ClockTool) Schema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"format": map[string]interface{}{
				"type":        "string",
				"description": "Optional format string (e.g. '2006-01-02 15:04:05'). Default is standard format.",
			},
		},
	}
}

func (t *ClockTool) Execute(args json.RawMessage) (string, error) {
	return time.Now().Format("2006-01-02 15:04:05 Monday"), nil
}
