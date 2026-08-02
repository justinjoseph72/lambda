### Go lang lambda Example
This is a project which will create an AWS lambda function which will listen on events AWS api gateway and will process a GET and POST request.

### Pre Requisite
Linux or Mac Os
Go 1.21 ++
AWS account and user with sufficient privilages to create cloud formation stack and roles
aws cli 2.34.9


### Stack
The project uses cloudformation template which will create the following resources
* An AWS lambda function based on Go lang. The handler is stored in an S3 bucket which is passed as a paramter to the cloud formation template
* Roles to create cloud watch logs for the lambda
* Api gateway with a GET and POST endpoint
* A production stage for the api gateway
* Api gateway key
* Cloudwatch logs for the gateway
* Dynamo db creation
* Invocation of the lambda function on the GET and POST method of the API and interact with S3 bucket and dynamodb

### Deployment 
#### Build Lambda function
The project uses a Makefile which abstracts the commands required to run and deploy the application in AWS.
The Lambda function uses Amazon linux 2023 runtime and runs the binary file created by the go code directly as handler.
The binary is packaged as a zip file and saved to an S3 bucket. The S3 bucket name is set in the environment variable `SOURCE_BUCKET` and is used by the Make file directly.
The lambda also reads from an S3 bucket and copies files to a destination bucket.
The source document bucket is read from `TEMP_DOC_BUCKET` environment variable and is copied to S3 bucket set in `FINAL_DOC_BUCKET` evironment variable.

```
export SOURCE_BUCKET=<the bucket name>
```
The lambda function can be build using the following command.

```
make build
```
This will build the binary executable with name `bootstrap` and it will ziped into a file set in the variable `zipFileName`. 

The binary can be build using the following command

```
make build zipFileName=<functionName.zip>
```

#### Upload lambda function to S3 bucket 
The following command will build and then upload the function to the S3 bucket defined in the environment variable `SOURCE_BUCKET`

```
make upload-lambda zipFile=<yourFunctionName.zip>
```

#### Creating the stack
The following command will create the stack using the cloudformation template defined in the variable `templateFile`

```
make createStack
```

#### Check status
Find the status of the stack using the command
```
make stackStatus
```

####
Once the stack is successfully created the invoke curl can be done by the following command

```
make findUrlDetails
```

#### Updating the lambda with new version and deploying the stack
The lambda code can be changed and updated to the current stack
We need to first delete the current zip file in the bucket

```
make remove-lambda
```

The build and upload the new version of the lambda 
```
make upload-lambda
```

Then deploy the lambda
```
make deployt-updated-lambda
```

#### Check status of stack
```
make stackStatus
```

### List events of stack
```
make stackEvents
```

#### Deleting the stack

```
make deleteStack
```


### Api

The lambda function get trigged for the following api invocation

* GET /prod/dummy. All the orders in the dummy-order dynamo db table will be returned.
* GET /prod/dummy?orderId=T343431. Items for specific order item will be returned
* GET /prod/dummy?fileKey=T343431. The file with key T34341 in TEMP_DOC_BUCKET will be copied to FINAL_DOC_BUCKET in copy folder
* POST /prod/dummy will save the order provided in the json
```
{
    "orderId": "T34344",
    "amount": 99.99,
    "item": "Gadget"
}
```

### Dynamodb stream processing
The dynamodb table will not stream the changes. This will be consumed by the lambda in file scripts/list_aws_resources.py
There is a lambda function ProcessEventLambda created for this which has roles to act on the stream data and is mapped to the stream using event source mapping EventSourceDummyTableStream


