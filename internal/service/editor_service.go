package service

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tranducquang/kiroku/internal/config"
)

// EditorService handles external editor integration
type EditorService struct {
	cfg *config.Config
}

// NewEditorService creates a new editor service
func NewEditorService(cfg *config.Config) *EditorService {
	return &EditorService{cfg: cfg}
}

// EditNote opens a note in the external editor and returns the updated content
func (s *EditorService) EditNote(title, content string) (newTitle, newContent string, err error) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "kiroku-*.md")
	if err != nil {
		return "", "", fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write content with title as first line
	fullContent := fmt.Sprintf("# %s\n\n%s", title, content)
	if _, err := tmpFile.WriteString(fullContent); err != nil {
		return "", "", fmt.Errorf("write temp file: %w", err)
	}
	tmpFile.Close()

	// Open editor
	cmd := exec.Command(s.cfg.Editor.Command, append(s.cfg.Editor.Args, tmpFile.Name())...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("run editor: %w", err)
	}

	// Read updated content
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", "", fmt.Errorf("read temp file: %w", err)
	}

	// Parse title and content
	lines := strings.SplitN(string(data), "\n", 3)
	if len(lines) == 0 {
		return title, content, nil
	}

	// Extract title from first line (remove # prefix)
	newTitle = strings.TrimPrefix(lines[0], "# ")
	newTitle = strings.TrimSpace(newTitle)
	if newTitle == "" {
		newTitle = title
	}

	// Rest is content
	if len(lines) > 2 {
		newContent = strings.TrimSpace(lines[2])
	} else if len(lines) > 1 {
		newContent = strings.TrimSpace(lines[1])
	}

	return newTitle, newContent, nil
}

// CreateNote opens an editor for a new note
func (s *EditorService) CreateNote(templateContent string) (title, content string, err error) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "kiroku-*.md")
	if err != nil {
		return "", "", fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write template content
	if templateContent != "" {
		if _, err := tmpFile.WriteString(templateContent); err != nil {
			return "", "", fmt.Errorf("write temp file: %w", err)
		}
	} else {
		if _, err := tmpFile.WriteString("# \n\n"); err != nil {
			return "", "", fmt.Errorf("write temp file: %w", err)
		}
	}
	tmpFile.Close()

	// Open editor
	cmd := exec.Command(s.cfg.Editor.Command, append(s.cfg.Editor.Args, tmpFile.Name())...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("run editor: %w", err)
	}

	// Read content
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", "", fmt.Errorf("read temp file: %w", err)
	}

	// Parse title and content
	lines := strings.SplitN(string(data), "\n", 3)
	if len(lines) == 0 {
		return "", "", nil
	}

	// Extract title from first line
	title = strings.TrimPrefix(lines[0], "# ")
	title = strings.TrimSpace(title)

	// Rest is content
	if len(lines) > 2 {
		content = strings.TrimSpace(lines[2])
	} else if len(lines) > 1 {
		content = strings.TrimSpace(lines[1])
	}

	return title, content, nil
}
