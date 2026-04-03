### Issues 

#### Listing of objects in Bucket
 * The create stack works and the lmabda execution role has the policies however the lambda function fails to list to objects in the S3 bucket. :white_check_mark: [Fixed with policy separate for bucket and bucket objects in the lambda exection role]
 * The lambda interact with the S3 bucket and DynamoDb in an insecure fashion. Need to build VPC endpoints for the lambda to interact with them and not go over the internet.