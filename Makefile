zipFileName := "myFunction10.zip"
executableName := "bootstrap"
stackName := "go-lambda-api-stack"
# templateFile := "file://template.yaml"
templateFile := "file://template-cw-loggin.yaml"

build:
	@echo " running target $@ Building Go executable for Linux AMD64"
	GOOS=linux GOARCH=amd64 go build -o $(executableName) cmd/main.go
	@echo "Zipping executable into $(zipFileName)"
	zip -r $(zipFileName) . -x "*.md" -x ".git/*" -x ".gitignore" -x ".DS_Store" -x ".aws-sam/*" -x "Makefile" -x "buildCurl.sh" -x "template.yaml" -x "template-cw-loggin.yaml" -x ".claude/*"
clean:
	@echo "Cleaning up..."
	rm $(executableName)
	rm *.zip

upload-lambda: build
	aws s3 cp $(zipFileName) s3://$(SOURCE_BUCKET)/$(zipFileName)	

listS3:
	echo ${SOURCE_BUCKET}
	aws s3 ls s3://$(SOURCE_BUCKET)

remove-lambda:
	aws s3 rm s3://$(SOURCE_BUCKET)/$(zipFileName)	

createStack:
	aws cloudformation create-stack \
  --stack-name $(stackName) \
  --template-body $(templateFile) \
  --parameters \
    ParameterKey=LambdaS3Bucket,ParameterValue=${SOURCE_BUCKET} \
    ParameterKey=LambdaS3Key,ParameterValue=${zipFileName} \
	ParameterKey=TempDocS3Bucket,ParameterValue=${TEMP_DOC_BUCKET} \
	ParameterKey=FinalDocS3Bucket,ParameterValue=${FINAL_DOC_BUCKET} \
  --capabilities CAPABILITY_IAM

deploy-updated-lambda: upload-lambda
	aws cloudformation update-stack \
  --stack-name $(stackName) \
  --use-previous-template \
  --parameters \
	ParameterKey=LambdaS3Bucket,ParameterValue=${SOURCE_BUCKET} \
	ParameterKey=LambdaS3Key,ParameterValue=${zipFileName} \
  --capabilities CAPABILITY_IAM 

deleteStack:
	aws cloudformation delete-stack \
  --stack-name $(stackName) \
  --region eu-west-2

findUrlDetails:
	./buildCurl.sh ${stackName} 

stackStatus:
	aws cloudformation describe-stacks --query "Stacks[][StackName,StackStatus,StackId,CreationTime]" --output text

stackEvents:
	aws cloudformation describe-stack-events --stack-name $(stackName) --query "StackEvents[][Timestamp,ResourceStatus,ResourceType,LogicalResourceId,ResourceStatusReason]" --output text

