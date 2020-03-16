#!/bin/bash

BUCKET_PREFIX=swoldemi
APPLICATION=amazon-cognito-domain-distribution
ACCOUNT=273450712882
declare -a REGIONS=(
    "us-east-1"
    "us-east-2"
    "us-west-2"
    "ap-south-1"
    "ap-northeast-2"
    "ap-southeast-1"
    "ap-southeast-2"
    "ap-northeast-2"
    "ca-central-1"
    "eu-central-1"
    "eu-west-1"
    "eu-west-2"
)

publish_all_regions()
{   
    for REGION in "${REGIONS[@]}"
    do  
        echo Deploying to region $REGION
	    sam package --template-file template.yaml --s3-bucket $BUCKET_PREFIX-$REGION --output-template-file packaged.yaml
        sam publish --region $REGION --template packaged.yaml
    done

    aws serverlessrepo put-application-policy \
        --region us-east-1 \
        --application-id arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-domain-distribution \
        --statements Principals=*,Actions=Deploy
}

echo "On branch `basename $CODEBUILD_WEBHOOK_HEAD_REF`"
if [ "`basename $CODEBUILD_WEBHOOK_HEAD_REF`" = "master" ]
then 
    publish_all_regions
else
    echo Skipping publish
fi
