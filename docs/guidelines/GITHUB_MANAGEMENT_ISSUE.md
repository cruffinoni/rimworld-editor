# GitHub Issue Management Guidelines for Therapôn Project

When creating or managing GitHub issues for this project:

## Style Requirements
- Keep descriptions concise and focused on the problem/need
- Remove all code examples and implementation details from issues
- Use clear, actionable language without excessive formatting
- Remove emojis and verbose explanations
- Focus on WHAT needs to be done, not HOW to do it

## Available Labels
- `bug` - Something isn't working
- `enhancement` - New feature or request
- `documentation` - Improvements or additions to docs
- `severity: critical/high/medium/low` - Priority levels

## Issue Structure
```
## Problem
[Clear problem statement]

## Requirements/What's Needed
[Bullet points of requirements]

## Impact
[Brief impact description]
```

## Project Context
- Therapôn is a Go-based event-driven chatbot framework
- Has LLM providers: Anthropic, Google, OpenAI, OpenAI-compatible
- Uses capabilities/tools system for LLM function calling
- Has MCP (Model Context Protocol) integration
- Built with event bus architecture

## Key Commands
- `gh issue create --title "Title" --label "enhancement" --body "$(cat <<'EOF' ... EOF)"`
- `gh issue edit NUMBER --body "$(cat <<'EOF' ... EOF)"`
- `gh issue comment NUMBER --body "text"`
- `gh issue list --state open`

## Remember
Keep issues as problem statements and requirements, not implementation guides.