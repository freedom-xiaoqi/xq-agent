package channels

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	webview "github.com/webview/webview_go"
)

//go:embed ui/index.html
var uiContent string

type WebviewChannel struct {
	w       webview.WebView
	handler func(Message)
	mu      sync.Mutex
	ready   bool
}

func NewWebviewChannel() *WebviewChannel {
	// Initialize webview
	// debug=true for development
	w := webview.New(true)
	w.SetTitle("Desktop Agent")
	w.SetSize(800, 600, webview.HintNone)

	c := &WebviewChannel{
		w: w,
	}

	// Bind Go function to JS
	w.Bind("sendMessageToAgent", func(content string) {
		if c.handler != nil {
			c.handler(Message{
				ID:      fmt.Sprintf("webview-%d", time.Now().UnixNano()),
				Content: content,
				Sender:  "user",
				Channel: "webview",
			})
		}
	})

	// Set initial content
	w.SetHtml(uiContent)

	return c
}

func (c *WebviewChannel) Name() string { return "webview" }

func (c *WebviewChannel) Start() error {
	// Webview Run() blocks the main thread, so it must be called from main()
	// This Start() method is just a placeholder or for initialization if needed.
	// We will manage the Run loop in main.go
	c.ready = true
	return nil
}

func (c *WebviewChannel) Run() {
	c.w.Run()
}

func (c *WebviewChannel) Stop() error {
	c.w.Terminate()
	return nil
}

func (c *WebviewChannel) SendMessage(content string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Escape content for JS string
	jsContent, _ := json.Marshal(content)
	script := fmt.Sprintf("window.appendMessage('agent', %s)", string(jsContent))

	// Ensure we are on the main thread or dispatch properly
	c.w.Dispatch(func() {
		c.w.Eval(script)
	})
	return nil
}

func (c *WebviewChannel) SendToken(token string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	jsToken, _ := json.Marshal(token)
	script := fmt.Sprintf("window.appendToken(%s)", string(jsToken))

	c.w.Dispatch(func() {
		c.w.Eval(script)
	})
	return nil
}

func (c *WebviewChannel) ShowThinking() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	script := "window.showThinking()"

	c.w.Dispatch(func() {
		c.w.Eval(script)
	})
	return nil
}

func (c *WebviewChannel) SendReasoning(content string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	jsContent, _ := json.Marshal(content)
	script := fmt.Sprintf("window.appendReasoning(%s)", string(jsContent))

	c.w.Dispatch(func() {
		c.w.Eval(script)
	})
	return nil
}

func (c *WebviewChannel) SendToolCall(toolName, args string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	jsName, _ := json.Marshal(toolName)
	jsArgs, _ := json.Marshal(args)
	script := fmt.Sprintf("window.appendToolCall(%s, %s)", string(jsName), string(jsArgs))

	c.w.Dispatch(func() {
		c.w.Eval(script)
	})
	return nil
}

func (c *WebviewChannel) IsStreamable() bool {
	return true
}

func (c *WebviewChannel) OnMessage(handler func(Message)) {
	c.handler = handler
}
