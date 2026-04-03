---
name: commit-code
# prettier-ignore
description: A skill to check the code before commiting to ensure there are no hardcoded values, AWS credentials, and that the README.md file is updated with any new instructions or changes to the code.
---

When commiting code:
- make sure there are no AWS credentials in the code
- make sure there are no hardcoded values that should be in environment variables or configuration files
- make sure to update the README.md file with any new instructions or changes to the code
- make sure to update the .gitignore file if any claude related files are added
- make sure no AWS account numbers are hardcoded especailly in and .md files and ARN references
- make sure no api key is exposed in the code or in any .md files
- The commit message should have summary of the changes made and reference any relevant issues or tickets. It should also follow the standard commit message format, which includes a type (e.g. feat, fix, docs), a scope (optional), and a concise description of the change. For example: "feat: add new API endpoint for retrieving orders".