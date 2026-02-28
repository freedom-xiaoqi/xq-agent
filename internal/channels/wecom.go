package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"xq-agent/internal/config"
)

type WeComChannel struct {
	cfg          config.WeComConfig
	handler      func(Message)
	accessToken  string
	tokenExpires time.Time
	mu           sync.Mutex
}

func NewWeComChannel(cfg config.WeComConfig) *WeComChannel {
	return &WeComChannel{
		cfg: cfg,
	}
}

func (c *WeComChannel) Name() string { return "wecom" }

func (c *WeComChannel) Start() error {
	if !c.cfg.Enabled {
		return nil
	}
	// In a real scenario, you would start an HTTP server here to receive callbacks
	// For now, we just ensure we can get a token
	if _, err := c.getAccessToken(); err != nil {
		fmt.Printf("[WeCom] Warning: Failed to get access token on start: %v\n", err)
	} else {
		fmt.Println("[WeCom] Successfully connected (AccessToken acquired).")
	}
	return nil
}

func (c *WeComChannel) Stop() error {
	return nil
}

func (c *WeComChannel) SendMessage(content string) error {
	if !c.cfg.Enabled {
		return nil
	}

	token, err := c.getAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %v", err)
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token)

	payload := map[string]interface{}{
		"touser":  "@all", // Send to all users in the app scope by default
		"msgtype": "text",
		"agentid": c.cfg.AgentID,
		"text": map[string]string{
			"content": content,
		},
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return err
	}

	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		return fmt.Errorf("wecom api error: %v, msg: %v", errcode, result["errmsg"])
	}

	fmt.Printf("[WeCom] Sent: %s\n", content)
	return nil
}

func (c *WeComChannel) SendToken(token string) error {
	// Not supported yet
	return nil
}

func (c *WeComChannel) ShowThinking() error {
	// WeCom doesn't support a "thinking" state in the same way, or we could send a temporary message.
	// For now, do nothing.
	return nil
}

func (c *WeComChannel) SendReasoning(content string) error {
	return nil
}

func (c *WeComChannel) SendToolCall(toolName, args string) error {
	// Maybe send a separate message?
	// c.SendMessage(fmt.Sprintf("[Thinking: Using Tool %s]", toolName))
	return nil
}

func (c *WeComChannel) IsStreamable() bool {
	return false
}

func (c *WeComChannel) OnMessage(handler func(Message)) {
	c.handler = handler
}

func (c *WeComChannel) getAccessToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.accessToken != "" && time.Now().Before(c.tokenExpires) {
		return c.accessToken, nil
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", c.cfg.CorpID, c.cfg.Secret)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("api error: %s", result.ErrMsg)
	}

	c.accessToken = result.AccessToken
	// Refresh a bit earlier than actual expiration
	c.tokenExpires = time.Now().Add(time.Duration(result.ExpiresIn-200) * time.Second)
	return c.accessToken, nil
}
