#!/bin/bash

PROJECT=go-dynamodb

aws cloudformation create-stack --stack-name go-dynamodb \
  --template-body file://main.yml \
  # --parameters ParameterKey=AssetsBucket,ParameterValue=$ASSETS_BUCKET \
  # ParameterKey=ExternalId,ParameterValue=$EXTERNAL_ID \
  # ParameterKey=ProjectName,ParameterValue=$PROJECT \
  # ParameterKey=Repository,ParameterValue=$REPOSITORY \
  --tags Key=project,Value=$PROJECT \
  --region ap-southeast-2 \
  --capabilities CAPABILITY_NAMED_IAM \
  --output yaml \
  --profile $PROFILE
