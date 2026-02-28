package tools

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileReadTool struct{}

func (t *FileReadTool) Name() string        { return "file_read" }
func (t *FileReadTool) Description() string { return "Read the contents of a file." }
func (t *FileReadTool) Schema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "The path to the file",
			},
		},
		"required": []string{"path"},
	}
}
func (t *FileReadTool) Execute(args json.RawMessage) (string, error) {
	var input struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return "", fmt.Errorf("invalid arguments: %v", err)
	}
	content, err := os.ReadFile(input.Path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

type FileListTool struct{}

func (t *FileListTool) Name() string        { return "file_list" }
func (t *FileListTool) Description() string { return "List files in a directory." }
func (t *FileListTool) Schema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "The directory path",
			},
		},
		"required": []string{"path"},
	}
}
func (t *FileListTool) Execute(args json.RawMessage) (string, error) {
	var input struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return "", fmt.Errorf("invalid arguments: %v", err)
	}
	entries, err := os.ReadDir(input.Path)
	if err != nil {
		return "", err
	}
	var result string
	for _, entry := range entries {
		info, _ := entry.Info()
		result += fmt.Sprintf("%s (%d bytes)\n", entry.Name(), info.Size())
	}
	return result, nil
}

type FileWriteTool struct{}

func (t *FileWriteTool) Name() string        { return "file_write" }
func (t *FileWriteTool) Description() string { return "Write content to a file (overwrite)." }
func (t *FileWriteTool) Schema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "The path to the file",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "The content to write",
			},
		},
		"required": []string{"path", "content"},
	}
}
func (t *FileWriteTool) Execute(args json.RawMessage) (string, error) {
	var input struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return "", fmt.Errorf("invalid arguments: %v", err)
	}
	err := os.WriteFile(input.Path, []byte(input.Content), 0644)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("File written to %s", input.Path), nil
}
