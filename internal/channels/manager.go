package channels

type Manager struct {
	channels []Channel
	msgChan  chan Message
}

func NewManager() *Manager {
	return &Manager{
		msgChan: make(chan Message, 100),
	}
}

func (m *Manager) Register(c Channel) {
	m.channels = append(m.channels, c)
	c.OnMessage(func(msg Message) {
		m.msgChan <- msg
	})
}

func (m *Manager) Start() {
	for _, c := range m.channels {
		c.Start()
	}
}

func (m *Manager) Stop() {
	for _, c := range m.channels {
		c.Stop()
	}
}

func (m *Manager) Messages() <-chan Message {
	return m.msgChan
}

func (m *Manager) Broadcast(content string) {
	for _, c := range m.channels {
		c.SendMessage(content)
	}
}

func (m *Manager) SendToChannel(channelName, content string) {
	for _, c := range m.channels {
		if c.Name() == channelName {
			c.SendMessage(content)
		}
	}
}

func (m *Manager) SendTokenToChannel(channelName, token string) {
	for _, c := range m.channels {
		if c.Name() == channelName {
			c.SendToken(token)
		}
	}
}

func (m *Manager) SendReasoningToChannel(channelName, content string) {
	for _, c := range m.channels {
		if c.Name() == channelName {
			c.SendReasoning(content)
		}
	}
}

func (m *Manager) SendToolCallToChannel(channelName, toolName, args string) {
	for _, c := range m.channels {
		if c.Name() == channelName {
			c.SendToolCall(toolName, args)
		}
	}
}

func (m *Manager) ShowThinking(channelName string) {
	for _, c := range m.channels {
		if c.Name() == channelName {
			c.ShowThinking()
		}
	}
}

func (m *Manager) IsChannelStreamable(channelName string) bool {
	for _, c := range m.channels {
		if c.Name() == channelName {
			return c.IsStreamable()
		}
	}
	return false
}

// InjectMessage allows internal components to send messages as if they came from a channel
func (m *Manager) InjectMessage(msg Message) {
	m.msgChan <- msg
}
