package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type BrowserOpenTool struct{}

func (t *BrowserOpenTool) Name() string { return "browser_open" }
func (t *BrowserOpenTool) Description() string {
	return "Open a URL in a headless browser and return the text content of the page."
}
func (t *BrowserOpenTool) Schema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"url": map[string]interface{}{
				"type":        "string",
				"description": "The URL to visit",
			},
		},
		"required": []string{"url"},
	}
}
func (t *BrowserOpenTool) Execute(args json.RawMessage) (string, error) {
	var input struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return "", fmt.Errorf("invalid arguments: %v", err)
	}
	input.URL = strings.TrimSpace(input.URL)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-software-rasterizer", true),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(input.URL),
		chromedp.Text("body", &res),
	)
	if err != nil {
		return "", fmt.Errorf("browser error: %v", err)
	}
	if len(res) > 2000 {
		res = res[:2000] + "...(truncated)"
	}
	return res, nil
}

type BrowserScreenshotTool struct{}

func (t *BrowserScreenshotTool) Name() string { return "browser_screenshot" }
func (t *BrowserScreenshotTool) Description() string {
	return "Take a screenshot of a URL and save it to a file."
}
func (t *BrowserScreenshotTool) Schema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"url": map[string]interface{}{
				"type":        "string",
				"description": "The URL to visit",
			},
			"output": map[string]interface{}{
				"type":        "string",
				"description": "Output file path (e.g. screenshot.png)",
			},
		},
		"required": []string{"url", "output"},
	}
}
func (t *BrowserScreenshotTool) Execute(args json.RawMessage) (string, error) {
	var input struct {
		URL    string `json:"url"`
		Output string `json:"output"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return "", err
	}
	input.URL = strings.TrimSpace(input.URL)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-software-rasterizer", true),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(input.URL),
		chromedp.CaptureScreenshot(&buf),
	)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(input.Output, buf, 0644); err != nil {
		return "", fmt.Errorf("failed to save screenshot: %v", err)
	}
	return fmt.Sprintf("Screenshot saved to %s", input.Output), nil
}
