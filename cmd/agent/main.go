package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/myuser/go-agent/internal/channels"
	"github.com/myuser/go-agent/internal/config"
	"github.com/myuser/go-agent/internal/core"
	"github.com/myuser/go-agent/internal/cron"
	"github.com/myuser/go-agent/internal/llm"
	"github.com/myuser/go-agent/internal/skills"
	"github.com/myuser/go-agent/internal/tools"
)

func main() {
	configFile := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// Load .env if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load config
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize LLM
	llmProvider := llm.NewOpenAI(cfg.LLM)

	// Initialize Channels
	cm := channels.NewManager()

	// Add Console Channel
	consoleCh := channels.NewConsoleChannel()
	cm.Register(consoleCh)

	// Add Webview Channel (GUI)
	// We always register it, but we need to handle its run loop specially
	webviewCh := channels.NewWebviewChannel()
	cm.Register(webviewCh)

	if cfg.Channels.Telegram.Enabled {
		cm.Register(channels.NewTelegramChannel(cfg.Channels.Telegram))
	}
	if cfg.Channels.WeCom.Enabled {
		cm.Register(channels.NewWeComChannel(cfg.Channels.WeCom))
	}

	// Initialize Skills Manager
	sm := skills.NewManager("skills")
	if err := sm.Load(); err != nil {
		log.Printf("Failed to load skills: %v", err)
	}

	// Initialize Cron Manager
	cronMgr := cron.NewManager(cm)
	cronMgr.Start()
	defer cronMgr.Stop()

	// Initialize Agent
	agent := core.NewAgent(cfg, llmProvider, cm)

	// Register Native Tools
	if cfg.Tools.BrowserEnabled {
		agent.RegisterTool(&tools.BrowserOpenTool{})
		agent.RegisterTool(&tools.BrowserScreenshotTool{})
	}
	if cfg.Tools.FileEnabled {
		agent.RegisterTool(&tools.FileReadTool{})
		agent.RegisterTool(&tools.FileListTool{})
		agent.RegisterTool(&tools.FileWriteTool{})
	}
	if cfg.Tools.ShellEnabled {
		agent.RegisterTool(&tools.ShellRunTool{})
	}
	// Always register Clock Tool (useful for cron jobs and time checks)
	agent.RegisterTool(&tools.ClockTool{})

	// Register Cron Tools
	for _, t := range cronMgr.Tools() {
		agent.RegisterTool(t.(tools.Tool))
	}

	// Inject Skills Context
	// We need to pass the skill descriptions to the Agent, either as tools or system prompt
	// Since we decided to let LLM read descriptions and use CLI, we append to system prompt.
	// We need to modify Agent to accept system prompt or modify history.
	// For now, let's inject a system message at startup.
	skillContext := sm.GetContext()
	if skillContext != "" {
		// Hack: Inject as a system message into the agent's history
		// Since agent.history is private, we should expose a method or pass it in constructor.
		// Let's modify Agent to accept initial system prompt.
		agent.SetSystemPrompt("You are a helpful AI agent. " + skillContext)
	}

	// Start Agent in a goroutine because Webview needs the main thread
	log.Println("Starting agent...")
	go func() {
		agent.Run()
	}()

	// Wait for interrupt in a separate goroutine if we are running webview
	// But webview.Run() blocks.
	// So we need to handle shutdown gracefully.

	// Start Webview (Blocks here)
	log.Println("Starting Webview GUI...")
	webviewCh.Run()

	// When webview closes, we exit
	log.Println("Shutting down...")
	cm.Stop()
}
