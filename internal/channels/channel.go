package channels

type Message struct {
	ID      string
	Content string
	Sender  string
	Channel string // e.g. "wecom", "dingtalk"
}

type Channel interface {
	Name() string
	Start() error
	Stop() error
	SendMessage(content string) error
	SendToken(token string) error
	ShowThinking() error
	SendReasoning(content string) error
	SendToolCall(toolName, args string) error
	IsStreamable() bool
	OnMessage(handler func(Message))
}
