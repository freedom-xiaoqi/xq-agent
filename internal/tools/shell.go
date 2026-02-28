package tools

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

type ShellRunTool struct{}

func (t *ShellRunTool) Name() string { return "shell_run" }
func (t *ShellRunTool) Description() string {
	return "Run a shell command and return the output. Use with caution."
}

func (t *ShellRunTool) Schema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": map[string]interface{}{
				"type":        "string",
				"description": "The command to run",
			},
		},
		"required": []string{"command"},
	}
}

func (t *ShellRunTool) Execute(args json.RawMessage) (string, error) {
	var input struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return "", fmt.Errorf("invalid arguments: %v", err)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command", input.Command)
	} else {
		cmd = exec.Command("bash", "-c", input.Command)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command failed: %v, output: %s", err, output)
	}

	return string(output), nil
}
