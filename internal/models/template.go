package models

import (
	"encoding/json"
	"errors"
	"time"
)

// Template types
const (
	TemplateTypeNote = "note"
	TemplateTypeTodo = "todo"
)

// ErrEmptyTemplateName is returned when template name is empty
var ErrEmptyTemplateName = errors.New("template name cannot be empty")

// TemplateVariable represents a template variable
type TemplateVariable struct {
	Name    string `json:"name"`
	Default string `json:"default"`
}

// Template represents a note template
type Template struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Type        string    `json:"type"`
	Icon        string    `json:"icon"`
	Variables   string    `json:"variables"` // JSON string
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Parsed variables (not stored in DB)
	ParsedVariables []TemplateVariable `json:"-"`
}

// Validate validates the template fields
func (t *Template) Validate() error {
	if t.Name == "" {
		return ErrEmptyTemplateName
	}
	if t.Type != TemplateTypeNote && t.Type != TemplateTypeTodo {
		t.Type = TemplateTypeNote
	}
	if t.Icon == "" {
		if t.Type == TemplateTypeTodo {
			t.Icon = "☐"
		} else {
			t.Icon = "📄"
		}
	}
	return nil
}

// GetVariables parses and returns template variables
func (t *Template) GetVariables() ([]TemplateVariable, error) {
	if t.Variables == "" {
		return nil, nil
	}

	var vars []TemplateVariable
	if err := json.Unmarshal([]byte(t.Variables), &vars); err != nil {
		return nil, err
	}

	t.ParsedVariables = vars
	return vars, nil
}

// SetVariables sets the template variables as JSON
func (t *Template) SetVariables(vars []TemplateVariable) error {
	if len(vars) == 0 {
		t.Variables = ""
		t.ParsedVariables = nil
		return nil
	}

	data, err := json.Marshal(vars)
	if err != nil {
		return err
	}

	t.Variables = string(data)
	t.ParsedVariables = vars
	return nil
}

// DefaultTemplates returns a list of default templates
func DefaultTemplates() []Template {
	return []Template{
		{
			Name:        "Blank Note",
			Description: "Empty note",
			Content:     "",
			Type:        TemplateTypeNote,
			Icon:        "📝",
			IsDefault:   true,
		},
		{
			Name:        "Blank Todo",
			Description: "Empty todo item",
			Content:     "",
			Type:        TemplateTypeTodo,
			Icon:        "☐",
			IsDefault:   true,
		},
		{
			Name:        "Meeting Notes",
			Description: "Template for meeting notes",
			Content: `# {{title}}

**Date:** {{date}}
**Attendees:** 

## Agenda
- 

## Discussion Points


## Action Items
- [ ] 

## Next Steps
`,
			Type: TemplateTypeNote,
			Icon: "🤝",
		},
		{
			Name:        "Daily Standup",
			Description: "Daily standup template",
			Content: `# Standup - {{date}}

## Yesterday
- 

## Today
- 

## Blockers
- None
`,
			Type: TemplateTypeNote,
			Icon: "☀️",
		},
		{
			Name:        "Bug Report",
			Description: "Bug report template",
			Content: `## Description


## Steps to Reproduce
1. 
2. 
3. 

## Expected Behavior


## Actual Behavior


## Environment
- OS: 
- Version:
`,
			Type: TemplateTypeTodo,
			Icon: "🐛",
		},
		{
			Name:        "Feature Request",
			Description: "Feature request template",
			Content: `## Summary


## User Story
As a [user type], I want [goal] so that [benefit].

## Acceptance Criteria
- [ ] 
- [ ] 

## Technical Notes
`,
			Type: TemplateTypeTodo,
			Icon: "✨",
		},
		{
			Name:        "Weekly Review",
			Description: "Weekly review template",
			Content: `# Week {{week_number}} Review

## Accomplishments
- 

## Challenges
- 

## Learnings
- 

## Next Week Goals
- [ ] 
- [ ]
`,
			Type: TemplateTypeNote,
			Icon: "📅",
		},
		{
			Name:        "Project Note",
			Description: "Project documentation template",
			Content: `# Project: {{title}}

## Overview


## Goals
- 

## Timeline
| Phase | Start | End | Status |
|-------|-------|-----|--------|
|       |       |     |        |

## Resources
- 

## Notes
`,
			Type: TemplateTypeNote,
			Icon: "📊",
		},
	}
}
