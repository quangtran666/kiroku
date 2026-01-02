-- migrations/002_seed.sql
-- Default data for Kiroku

-- Insert default folders
INSERT OR IGNORE INTO folders (id, name, icon) VALUES (1, 'Work', '💼');
INSERT OR IGNORE INTO folders (id, name, icon) VALUES (2, 'Personal', '🏠');
INSERT OR IGNORE INTO folders (id, name, icon) VALUES (3, 'Ideas', '💡');

-- Insert default templates
INSERT OR IGNORE INTO templates (name, description, content, type, icon, is_default) VALUES
    ('Blank Note', 'Empty note', '', 'note', '📝', TRUE);

INSERT OR IGNORE INTO templates (name, description, content, type, icon, is_default) VALUES
    ('Blank Todo', 'Empty todo item', '', 'todo', '☐', TRUE);

INSERT OR IGNORE INTO templates (name, description, content, type, icon, is_default) VALUES
    ('Meeting Notes', 'Template for meeting notes', '# {{title}}

**Date:** {{date}}
**Attendees:** 

## Agenda
- 

## Discussion Points


## Action Items
- [ ] 

## Next Steps
', 'note', '🤝', FALSE);

INSERT OR IGNORE INTO templates (name, description, content, type, icon, is_default) VALUES
    ('Daily Standup', 'Daily standup template', '# Standup - {{date}}

## Yesterday
- 

## Today
- 

## Blockers
- None
', 'note', '☀️', FALSE);

INSERT OR IGNORE INTO templates (name, description, content, type, icon, is_default) VALUES
    ('Bug Report', 'Bug report template', '## Description


## Steps to Reproduce
1. 
2. 
3. 

## Expected Behavior


## Actual Behavior


## Environment
- OS: 
- Version:
', 'todo', '🐛', FALSE);

INSERT OR IGNORE INTO templates (name, description, content, type, icon, is_default) VALUES
    ('Weekly Review', 'Weekly review template', '# Week {{week_number}} Review

## Accomplishments
- 

## Challenges
- 

## Learnings
- 

## Next Week Goals
- [ ] 
- [ ]
', 'note', '📅', FALSE);
