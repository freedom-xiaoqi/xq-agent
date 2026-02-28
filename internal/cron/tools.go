package cron

import (
	"encoding/json"
	"fmt"
)

type CronAddTool struct {
	manager *Manager
}

func (t *CronAddTool) Name() string { return "cron_add" }
func (t *CronAddTool) Description() string {
	return "Add a recurring task. Format: * * * * * * (seconds minutes hours day month weekday). Example: '0 30 8 * * *' for 8:30 AM daily."
}
func (t *CronAddTool) Schema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"spec": map[string]interface{}{
				"type":        "string",
				"description": "Cron expression (with seconds)",
			},
			"task": map[string]interface{}{
				"type":        "string",
				"description": "Description of the task to perform",
			},
		},
		"required": []string{"spec", "task"},
	}
}
func (t *CronAddTool) Execute(args json.RawMessage) (string, error) {
	var input struct {
		Spec string `json:"spec"`
		Task string `json:"task"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return "", err
	}
	id, err := t.manager.AddJob(input.Spec, input.Task)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Job added with ID %d", id), nil
}

type CronListTool struct {
	manager *Manager
}

func (t *CronListTool) Name() string        { return "cron_list" }
func (t *CronListTool) Description() string { return "List all active cron jobs." }
func (t *CronListTool) Schema() interface{} { return map[string]interface{}{"type": "object"} }
func (t *CronListTool) Execute(args json.RawMessage) (string, error) {
	jobs := t.manager.ListJobs()
	if len(jobs) == 0 {
		return "No active jobs.", nil
	}
	var res string
	for id, task := range jobs {
		res += fmt.Sprintf("ID %d: %s\n", id, task)
	}
	return res, nil
}

type CronRemoveTool struct {
	manager *Manager
}

func (t *CronRemoveTool) Name() string        { return "cron_remove" }
func (t *CronRemoveTool) Description() string { return "Remove a cron job by ID." }
func (t *CronRemoveTool) Schema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Job ID to remove",
			},
		},
		"required": []string{"id"},
	}
}
func (t *CronRemoveTool) Execute(args json.RawMessage) (string, error) {
	var input struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return "", err
	}
	t.manager.RemoveJob(input.ID)
	return fmt.Sprintf("Job %d removed", input.ID), nil
}

func (m *Manager) Tools() []interface{} {
	return []interface{}{
		&CronAddTool{manager: m},
		&CronListTool{manager: m},
		&CronRemoveTool{manager: m},
	}
}
