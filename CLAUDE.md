# CLAUDE.md

## Project Overview

A Go AWS Lambda learning project that implements a simple order-processing REST API exposed via API Gateway, with S3 integration. The Lambda handles GET and POST requests on the `/dummy` endpoint and lists objects from a configured S3 bucket on each invocation.

## Tech Stack

- **Language**: Go (compiled to `bootstrap` binary for `provided.al2023` runtime)
- **Target**: Linux AMD64 (`GOOS=linux GOARCH=amd64`)
- **AWS Services**: Lambda, API Gateway (REST), S3, CloudWatch Logs, CloudFormation, IAM
- **Key Dependencies**: `aws-lambda-go v1.52.0`, `aws-sdk-go-v2`

## Project Structure

```
cmd/main.go                       # Lambda handler entry point; initializes S3 client
internal/svc/process-event.go     # Event routing and response building
internal/util/utils.go            # JSON unmarshaling helpers
template.yaml                     # Basic CloudFormation template (not active)
template-cw-loggin.yaml           # Active CloudFormation template (includes CW logging, S3 IAM)
Makefile                          # All build and deployment automation
buildCurl.sh                      # Extracts API URL and generates curl test command
```

## Build & Deploy

All operations go through `make`. Key targets:

| Command | Action |
|---|---|
| `make build` | Compile Go → `bootstrap`, zip → `myFunction6.zip` |
| `make upload-lambda` | Upload zip to S3 (requires `SOURCE_BUCKET` env var) |
| `make createStack` | Deploy CloudFormation stack |
| `make deploy-updated-lambda` | Rebuild, re-upload, and update existing stack |
| `make stackStatus` | Check CloudFormation stack status |
| `make stackEvents` | View stack events |
| `make findUrlDetails` | Extract API Gateway URL |
| `make deleteStack` | Tear down all AWS resources |

**Required environment variable**: `SOURCE_BUCKET` (the S3 bucket used for Lambda code storage and runtime listing).

## CloudFormation

The active template is `template-cw-loggin.yaml`. It provisions:
- Lambda function (128MB, 10s timeout, `provided.al2023`)
- API Gateway REST API with `/dummy` GET and POST methods
- IAM execution role with S3 read/list/write permissions
- CloudWatch Log Group (14-day retention)
- API Gateway access logging to CloudWatch

Region is hardcoded to `eu-west-2` in Makefile delete targets.

## API Endpoints

- `GET /dummy` → Returns hardcoded sample order (`orderId: 12345, amount: 100.00`)
- `POST /dummy` → Parses order JSON from body, returns response with `createdAt` timestamp

## Known Issues / Notes

- S3 IAM policy requires **two separate statements**: one for the bucket itself (`arn:aws:s3:::bucket`) and one for objects (`arn:aws:s3:::bucket/*`). Using only one causes access errors.
- The Lambda also lists S3 objects on every invocation (logged output, not returned in API response).
- API Key is configured in CloudFormation but not enforced in Lambda code.
