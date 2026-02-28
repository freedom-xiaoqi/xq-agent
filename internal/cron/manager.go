package cron

import (
	"fmt"
	"sync"

	"github.com/myuser/go-agent/internal/channels"
	"github.com/robfig/cron/v3"
)

type Manager struct {
	cron *cron.Cron
	cm   *channels.Manager
	mu   sync.Mutex
	jobs map[cron.EntryID]string // ID -> Description
}

func NewManager(cm *channels.Manager) *Manager {
	return &Manager{
		cron: cron.New(cron.WithSeconds()), // Enable seconds field
		cm:   cm,
		jobs: make(map[cron.EntryID]string),
	}
}

func (m *Manager) Start() {
	m.cron.Start()
}

func (m *Manager) Stop() {
	m.cron.Stop()
}

func (m *Manager) AddJob(spec string, task string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id, err := m.cron.AddFunc(spec, func() {
		// Send a system message to trigger the agent
		taskName := task
		fmt.Printf("[CRON] Triggered: %s\n", taskName)

		// Direct injection via a special method we should add to Channel Manager
		m.cm.InjectMessage(channels.Message{
			ID:      fmt.Sprintf("cron-%s", taskName), // Use task name as part of ID
			Content: fmt.Sprintf("It is time to: %s", taskName),
			Sender:  "system_scheduler",
			Channel: "webview", // Default to webview for now
		})
	})

	if err != nil {
		return 0, err
	}
	m.jobs[id] = task
	return int(id), nil
}

func (m *Manager) ListJobs() map[int]string {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make(map[int]string)
	for id, desc := range m.jobs {
		result[int(id)] = desc
	}
	return result
}

func (m *Manager) RemoveJob(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cron.Remove(cron.EntryID(id))
	delete(m.jobs, cron.EntryID(id))
}
