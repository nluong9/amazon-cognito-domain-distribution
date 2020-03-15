## Examples of using amazon-cognito-custom-domain-link

These examples assume you are creating a [`AWS::Cognito::UserPoolDomain`](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cognito-userpooldomain.html) in the same stack that you are using this custom resource. Note that deleting and recreating a `AWS::Cognito::UserPoolDomain` can take 15 minutes to fully create, 20 minutes to delete, and 1 hour for the deletion to fully propagate through AWS if you are planning on attempting frequent creations and deletions. 

### 1. Serverless Application Model Template
It is recommended that you use this custom resource as a Severless Application Repository nested app in your CloudFormation template.
 - Example that will only return the DNS name of the CloudFront distribution: [here](./examples/sam/no-create-sam-template.yaml)
 - Example that also creates an alias record: [here](./examples/sam/sam-template.yaml)

### 2. CloudFormation Template
It is also possible to deploy this function seperately and link the Lambda function's ARN as the ServiceToken propery of the `AWS::CloudFormation::CustomResource`. The default values of the 3 parameters recognized by the SAR app are sufficient to safely deploy he application and retrieve the Lambda function's Amazon Resource Name. 
 - Example that will only return the DNS name of the CloudFront distribution: [here](./examples/cloudformation/no-create-cfn-template.yaml)
 - Example that also creates an alias record: [here](./examples/cloudformation/cfn-template.yaml)
