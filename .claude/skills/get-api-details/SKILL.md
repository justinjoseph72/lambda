---
name: get-api-details
# prettier-ignore
description: A skill to get the api details
disable-model-invocation: true
---

When asked about the API details, provide the following information:
- Run the script findApiUrl.sh to get the API URL and export it as an environment variable STACK_API_URL
- Run the script findApiKey.sh to get the API key value and export it as an environment variable STACK_API_KEY_VALUE
- Provide the curl command to test the API using the API URL and API key value. The command should be in the format: curl $STACK_API_URL -H "x-api-key: $STACK_API_KEY_VALUE"
