#!/usr/bin/env bash
set -e

stackName=$1

if [ -z "$stackName" ]; then
  echo "Usage: $0 <stack-name>"
  exit 1
fi


apiKeyId=$(aws cloudformation describe-stacks --stack-name $stackName \
--query "Stacks[0].Outputs[?OutputKey=='ApiKeyValue'].OutputValue" --output text)

apiKeyValue=$(aws apigateway get-api-key --api-key $apiKeyId --include-value | jq -r '.value')

export STACK_API_KEY_VALUE=$apiKeyValue

