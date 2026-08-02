import logging
import json
import boto3
from botocore.exceptions import ClientError

Logger = logging.getLogger()
Logger.setLevel(logging.INFO)
session = boto3.Session()
dynamodb = session.client('dynamodb')

def handler(event, context):
    Logger.info('Received Event: %s', event)
    for rec in event['Records']:
        Logger.info('Record: %s', rec)
        # if rec['eventName'] == 'INSERT':
        #     new_image = rec['dynamodb']['NewImage']
        #     order_id = new_image['orderId']['S']
        #     Logger.info('Processing new order with ID: %s', order_id)
        #     # Example: Update the order status in DynamoDB
        #     try:
        #         response = dynamodb.update_item(
        #             TableName='OrdersTable',
        #             Key={'orderId': {'S': order_id}},
        #             UpdateExpression='SET orderStatus = :status',
        #             ExpressionAttributeValues={':status': {'S': 'Processed'}}
        #         )
        #         Logger.info('Order %s status updated to Processed', order_id)
        #     except ClientError as e:
        #         Logger.error('Failed to update order %s: %s', order_id, e)
            