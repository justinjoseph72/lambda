#!/usr/bin/env bash
set -e

stackName=$1

if [ -z "$stackName" ]; then
  echo "Usage: $0 <stack-name>"
  exit 1
fi


apiUrl=$(aws cloudformation describe-stacks --stack-name $stackName \
--query "Stacks[0].Outputs[?OutputKey=='ApiInvokeURL'].OutputValue" --output text)

export STACK_API_URL=$apiUrl
echo $apiUrl
