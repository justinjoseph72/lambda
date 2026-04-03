---
description: Check code for security issues and hardcoded values, update README, then commit
---

Before committing, check the following and fix any issues found:

- No AWS credentials or API keys in any file (code, markdown, config)
- No hardcoded AWS account numbers, including in ARN references in .md files
- No hardcoded values that should be environment variables or config (e.g. bucket names, table names, region)
- No exposed API key values in code or markdown files
- README.md is updated with any new instructions or changes introduced by this commit
- .gitignore is updated if any new generated or sensitive files were added

Commit message format: `<type>(<scope>): <description>` (e.g. `feat: add DynamoDB integration for order persistence`). Include a short body summarising key changes if there is more than one logical change.
