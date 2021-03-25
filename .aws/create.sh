#!/bin/bash

PROJECT=go-dynamodb

aws cloudformation create-stack --stack-name go-dynamodb \
  --template-body file://main.yml \
  --parameters ParameterKey=TableName,ParameterValue=$TABLE_NAME \
  ParameterKey=ProjectName,ParameterValue=$PROJECT \
  --tags Key=project,Value=$PROJECT \
  --region ap-southeast-2 \
  --capabilities CAPABILITY_NAMED_IAM \
  --output yaml \
  --profile $PROFILE
