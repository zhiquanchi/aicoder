package tui

import (
	"fmt"
	"os/exec"
	"strings"
)

// ToolChecker provides methods to check tool installation status
type ToolChecker struct{}

func NewToolChecker() *ToolChecker {
	return &ToolChecker{}
}

// CheckToolStatus checks if a tool is installed and returns its status
func (tc *ToolChecker) CheckToolStatus(toolName string) (bool, string) {
	// Map tool names to their binary names
	binaryMap := map[string][]string{
		"Claude":     {"claude", "claude-code"},
		"Gemini":     {"gemini"},
		"Codex":      {"codex", "openai"},
		"OpenCode":   {"opencode"},
		"CodeBuddy":  {"codebuddy", "codebuddy-code"},
		"Qoder":      {"qodercli", "qoder"},
		"IFlow":      {"iflow"},
		"Kilo":       {"kilo"},
	}

	binaries, ok := binaryMap[toolName]
	if !ok {
		return false, "Unknown tool"
	}

	for _, binary := range binaries {
		if path, err := exec.LookPath(binary); err == nil {
			// Try to get version
			cmd := exec.Command(binary, "--version")
			output, err := cmd.CombinedOutput()
			version := "Installed"
			if err == nil {
				versionStr := strings.TrimSpace(string(output))
				if len(versionStr) > 0 && len(versionStr) < 100 {
					version = versionStr
				}
			}
			return true, fmt.Sprintf("✓ %s (%s)", path, version)
		}
	}

	return false, "✗ Not installed"
}

// GetAllToolStatuses returns the status of all supported tools
func (tc *ToolChecker) GetAllToolStatuses() map[string]string {
	tools := []string{"Claude", "Gemini", "Codex", "OpenCode", "CodeBuddy", "Qoder", "IFlow", "Kilo"}
	statuses := make(map[string]string)

	for _, tool := range tools {
		installed, status := tc.CheckToolStatus(tool)
		if installed {
			statuses[tool] = status
		} else {
			statuses[tool] = status
		}
	}

	return statuses
}
