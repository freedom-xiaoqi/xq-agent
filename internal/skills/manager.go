package skills

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Metadata struct {
	OpenClaw OpenClawMetadata `yaml:"openclaw"`
	Clawdbot ClawdbotMetadata `yaml:"clawdbot"` // Alias
}

type OpenClawMetadata struct {
	Requires map[string]interface{} `yaml:"requires"`
	Env      []string               `yaml:"env"`
	Bins     []string               `yaml:"bins"`
}

type ClawdbotMetadata struct {
	Requires map[string]interface{} `yaml:"requires"`
	Env      []string               `yaml:"env"`
	Bins     []string               `yaml:"bins"`
}

type Skill struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Metadata    Metadata `yaml:"metadata"`
	Path        string   `yaml:"-"` // Absolute path to skill directory
	Content     string   `yaml:"-"` // Full content of SKILL.md for LLM context
}

type Manager struct {
	skillsPath string
	skills     []Skill
}

func NewManager(path string) *Manager {
	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	// If path is relative, try to find it relative to executable
	if !filepath.IsAbs(path) {
		ex, err := os.Executable()
		if err == nil {
			exPath := filepath.Dir(ex)
			possiblePath := filepath.Join(exPath, path)
			if _, err := os.Stat(possiblePath); err == nil {
				absPath = possiblePath
			}
		}
	}

	return &Manager{
		skillsPath: absPath,
		skills:     make([]Skill, 0),
	}
}

func (m *Manager) Load() error {
	log.Printf("Loading skills from: %s", m.skillsPath)
	entries, err := os.ReadDir(m.skillsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No skills directory, that's fine
		}
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			skillPath := filepath.Join(m.skillsPath, entry.Name())
			skillFile := filepath.Join(skillPath, "SKILL.md")

			if _, err := os.Stat(skillFile); err == nil {
				skill, err := parseSkill(skillFile)
				if err != nil {
					fmt.Printf("Error parsing skill %s: %v\n", entry.Name(), err)
					continue
				}
				skill.Path = skillPath

				// Apply skill configuration
				applySkillConfig(skill)

				m.skills = append(m.skills, *skill)
				fmt.Printf("Loaded skill: %s\n", skill.Name)
			}
		}
	}
	return nil
}

func applySkillConfig(skill *Skill) {
	// 1. Set Environment Variables
	envVars := skill.Metadata.OpenClaw.Env
	if len(envVars) == 0 {
		envVars = skill.Metadata.Clawdbot.Env
	}
	for _, env := range envVars {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			os.Setenv(parts[0], parts[1])
			log.Printf("[Skill: %s] Set Env: %s", skill.Name, parts[0])
		}
	}

	// 2. Check Binary Dependencies
	bins := skill.Metadata.OpenClaw.Bins
	if len(bins) == 0 {
		bins = skill.Metadata.Clawdbot.Bins
	}
	for _, bin := range bins {
		if _, err := exec.LookPath(bin); err != nil {
			log.Printf("[Skill: %s] Warning: Required binary '%s' not found in PATH", skill.Name, bin)
		}
	}
}

func parseSkill(path string) (*Skill, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Extract Frontmatter
	sContent := string(content)
	if !strings.HasPrefix(sContent, "---") {
		return nil, fmt.Errorf("missing frontmatter")
	}

	parts := strings.SplitN(sContent, "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid frontmatter format")
	}

	var skill Skill
	if err := yaml.Unmarshal([]byte(parts[1]), &skill); err != nil {
		return nil, err
	}

	skill.Content = sContent // Keep full content for LLM context
	return &skill, nil
}

func (m *Manager) GetContext() string {
	var sb strings.Builder
	sb.WriteString("You have access to the following external skills (installed locally):\n\n")
	for _, s := range m.skills {
		sb.WriteString(fmt.Sprintf("## Skill: %s\n", s.Name))
		sb.WriteString(fmt.Sprintf("Description: %s\n", s.Description))
		sb.WriteString(fmt.Sprintf("Path: %s\n", s.Path))
		// We could include full content, but let's summarize or just put description to save tokens
		// If description is short, maybe include usage examples if parsed?
		// For now, let's assume the Description in frontmatter is enough, or append the whole SKILL.md if needed.
		// A better approach: "To use this skill, run commands as described below:"
		// Then append the body of SKILL.md (parts[2])

		// Let's include the full content for now as OpenClaw relies on it.
		// Truncate if too long?
		sb.WriteString("\n--- Skill Documentation ---\n")
		sb.WriteString(s.Content)
		sb.WriteString("\n---------------------------\n\n")
	}
	return sb.String()
}
