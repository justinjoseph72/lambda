---
description: Get the API URL and API key from the deployed CloudFormation stack and generate curl test commands
allowed-tools: Bash
---

Get the API details for the deployed stack by doing the following:

1. Run `.claude/skills/get-api-details/scripts/findApiUrl.sh go-lambda-api-stack` and capture the printed API URL
2. Run `.claude/skills/get-api-details/scripts/findApiKey.sh go-lambda-api-stack` and capture the printed API key value
3. Output:
   - The API URL
   - The API key value
   - A curl command for GET: `curl "<URL>" -H "x-api-key: <KEY>"`
   - A curl command for POST: `curl -X POST "<URL>" -H "x-api-key: <KEY>" -H "Content-Type: application/json" -d '{"orderId":1,"name":"Widget","amount":9.99}'`
