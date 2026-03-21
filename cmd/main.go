package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"

	svc "first-excercise/internal/svc"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client

func init() {
	log.Println("Lambda function initialized")
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
}

func handler(ctx context.Context, event json.RawMessage) (svc.EventResponse, error) {
	lc, _ := lambdacontext.FromContext(ctx)
	log.Printf("Lambda request ID: %s", lc.AwsRequestID)
	log.Printf("Received event: %s", string(event))
	log.Printf(" the arn is %s", lc.InvokedFunctionArn)
	bucket := os.Getenv("SOURCE_BUCKET_NAME")
	if bucket == "" {
		log.Printf("SOURCE_BUCKET_NAME environment variable is not set")
		return svc.BuildErrorResponse("SOURCE_BUCKET_NAME environment variable is not set"), fmt.Errorf("SOURCE_BUCKET_NAME environment variable is not set")
	}
	log.Printf("The configured bucket is %s \n", bucket)

	output, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucket,
	})
	if err != nil {
		log.Printf("Error listing objects in bucket: %v", err)
		return svc.BuildErrorResponse(fmt.Sprintf("Error listing objects in bucket: %v", err)), fmt.Errorf("Error listing objects in bucket: %v", err)
	}

	for _, object := range output.Contents {
		log.Printf("Object key: %s, Size: %d bytes\n", *object.Key, object.Size)
	}

	resultData, err := svc.ProcessEvent(event)
	if err != nil {
		cause := fmt.Sprintf("Error processing event: %v", err)
		return svc.BuildErrorResponse(cause), fmt.Errorf("failed to process event: %v", err)
	}

	if resultData.Success == true {
		log.Printf("succcesfully process the event")
	} else {
		log.Printf("failed to process the event with message: %s", resultData.Message)
	}
	log.Println("starting to build the response")
	resp := svc.BuildResponse(*resultData)

	return resp, nil

}

func main() {
	lambda.Start(handler)
}
