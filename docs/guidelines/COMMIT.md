# Commit Message Guidelines

Commit messages must follow the conventional commit format.

## Types

- **feat**: new feature for users
- **fix**: bug fix for users
- **docs**: documentation changes
- **style**: formatting, missing semicolons (no logic changes)
- **refactor**: code restructuring without behavior changes
- **test**: adding/updating tests (no production code changes)
- **chore**: build scripts, dependencies, configs (no production code changes)

## Format

**Header (required)**: `type(scope): description`
- Use imperative mood: "add", "fix", "update" (not "added", "fixes", "updating")
- Keep under 50 characters total
- No period at the end

**Body (optional)**:
- Leave blank line after header
- Explain what and why, not how
- Wrap at 72 characters per line
- Keep it short and concise

## Examples

```
feat: add user authentication

fix(api): resolve timeout on large requests

feat(llm): add Ollama embedding support

Integrates Ollama embeddings through OpenAI-compatible interface.
Enables local embedding generation for HAT operations.

chore: bump dependency versions
```

## General recommendations

- One logical change per commit
- Be specific but concise