package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LLM      LLMConfig      `yaml:"llm"`
	Channels ChannelsConfig `yaml:"channels"`
	Tools    ToolsConfig    `yaml:"tools"`
}

type LLMConfig struct {
	APIKey  string `yaml:"api_key"`
	BaseURL string `yaml:"base_url"`
	Model   string `yaml:"model"`
}

type ChannelsConfig struct {
	WeCom    WeComConfig    `yaml:"wecom"`
	DingTalk DingTalkConfig `yaml:"dingtalk"`
	Telegram TelegramConfig `yaml:"telegram"`
	Lark     LarkConfig     `yaml:"lark"`
}

type WeComConfig struct {
	Enabled bool   `yaml:"enabled"`
	CorpID  string `yaml:"corp_id"`
	AgentID int    `yaml:"agent_id"`
	Secret  string `yaml:"secret"`
}

type DingTalkConfig struct {
	Enabled   bool   `yaml:"enabled"`
	AppKey    string `yaml:"app_key"`
	AppSecret string `yaml:"app_secret"`
}

type TelegramConfig struct {
	Enabled bool   `yaml:"enabled"`
	Token   string `yaml:"token"`
}

type LarkConfig struct {
	Enabled   bool   `yaml:"enabled"`
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

type ToolsConfig struct {
	BrowserEnabled bool `yaml:"browser_enabled"`
	ShellEnabled   bool `yaml:"shell_enabled"`
	FileEnabled    bool `yaml:"file_enabled"`
	MCPEnabled     bool `yaml:"mcp_enabled"`
}

func Load(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
