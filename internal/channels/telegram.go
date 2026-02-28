package channels

import (
	"fmt"
	"time"

	"xq-agent/internal/config"
)

type TelegramChannel struct {
	cfg     config.TelegramConfig
	handler func(Message)
	stop    chan struct{}
}

func NewTelegramChannel(cfg config.TelegramConfig) *TelegramChannel {
	return &TelegramChannel{
		cfg:  cfg,
		stop: make(chan struct{}),
	}
}

func (c *TelegramChannel) Name() string { return "telegram" }

func (c *TelegramChannel) Start() error {
	if !c.cfg.Enabled {
		return nil
	}
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-c.stop:
				return
			case <-ticker.C:
				// Implement polling here
				// For now, it's a placeholder
			}
		}
	}()
	return nil
}

func (c *TelegramChannel) Stop() error {
	close(c.stop)
	return nil
}

func (c *TelegramChannel) SendMessage(content string) error {
	if !c.cfg.Enabled {
		return nil
	}
	fmt.Printf("[Telegram] Sending: %s\n", content)
	return nil
}

func (c *TelegramChannel) SendToken(token string) error {
	// Not supported yet
	return nil
}

func (c *TelegramChannel) ShowThinking() error {
	// Not supported yet
	return nil
}

func (c *TelegramChannel) SendReasoning(content string) error {
	// Not supported
	return nil
}

func (c *TelegramChannel) SendToolCall(toolName, args string) error {
	// Not supported
	return nil
}

func (c *TelegramChannel) IsStreamable() bool {
	return false
}

func (c *TelegramChannel) OnMessage(handler func(Message)) {
	c.handler = handler
}
