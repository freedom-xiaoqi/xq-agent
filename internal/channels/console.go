package channels

import (
	"bufio"
	"fmt"
	"os"
)

type ConsoleChannel struct {
	handler func(Message)
	stop    chan struct{}
}

func NewConsoleChannel() *ConsoleChannel {
	return &ConsoleChannel{
		stop: make(chan struct{}),
	}
}

func (c *ConsoleChannel) Name() string {
	return "console"
}

func (c *ConsoleChannel) Start() error {
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Console channel started. Type a message and press Enter.")
		for {
			select {
			case <-c.stop:
				return
			default:
				if scanner.Scan() {
					text := scanner.Text()
					if c.handler != nil {
						c.handler(Message{
							ID:      "console-msg",
							Content: text,
							Sender:  "user",
							Channel: "console",
						})
					}
				}
			}
		}
	}()
	return nil
}

func (c *ConsoleChannel) Stop() error {
	close(c.stop)
	return nil
}

func (c *ConsoleChannel) SendMessage(content string) error {
	fmt.Printf("[Agent]: %s\n", content)
	return nil
}

func (c *ConsoleChannel) SendToken(token string) error {
	fmt.Print(token)
	return nil
}

func (c *ConsoleChannel) ShowThinking() error {
	fmt.Print("Thinking...\r")
	return nil
}

func (c *ConsoleChannel) SendReasoning(content string) error {
	fmt.Print(content)
	return nil
}

func (c *ConsoleChannel) SendToolCall(toolName, args string) error {
	fmt.Printf("\n[Tool Call] %s(%s)\n", toolName, args)
	return nil
}

func (c *ConsoleChannel) IsStreamable() bool {
	return true
}

func (c *ConsoleChannel) OnMessage(handler func(Message)) {
	c.handler = handler
}
