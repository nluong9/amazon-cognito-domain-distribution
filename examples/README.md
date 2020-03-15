## Examples of using amazon-cognito-custom-domain-link

These examples assume you are creating a [`AWS::Cognito::UserPoolDomain`](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cognito-userpooldomain.html) in the same stack that you are using this custom resource. Note that deleting and recreating a `AWS::Cognito::UserPoolDomain` can take 15 minutes to fully create, 20 minutes to delete, and 1 hour for the deletion to fully propagate through AWS if you are planning on attempting frequent creations and deletions. It is also possible to deploy the Lambda function template seperately and interact with vanillia CloudFormation instead.

###  Serverless Application Model Template
It is recommended that you use this custom resource as a Severless Application Repository nested app.
 - Example that will only return the DNS name of the CloudFront distribution: [here](./examples/no-create-sam-template.yaml)
 - Example that also creates an alias record: [here](./examples/sam-template.yaml)
