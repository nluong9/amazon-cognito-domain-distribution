![](https://codebuild.us-east-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiU1NyMHI4KytFRzhZSUVEY2R0YTlwanBJTk9EdWNYbW93TzdRU3NCbUJ0TFZYMy9jUktROXlUQktEOUVjd0dJSDBWbHNtVjVqSFpaNWxvbTJxd0o4dW53PSIsIml2UGFyYW1ldGVyU3BlYyI6ImgyNlBtRXoyU1ZSNjNWZjYiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)
[![][sar-logo]](https://serverlessrepo.aws.amazon.com/applications/arn:aws:serverlessrepo:us-east-1:273450712882:applications~amazon-cognito-custom-domain-link)

[sar-deploy]: https://img.shields.io/badge/Serverless%20Application%20Repository-Deploy%20Now-FF9900?logo=amazon%20aws&style=flat-square
[sar-logo]: https://img.shields.io/badge/Serverless%20Application%20Repository-View-FF9900?logo=amazon%20aws&style=flat-square

# Amazon Cognito Custom Domain Link
>A CloudFormation Custom Resource for automatically linking your Cognito User Pool's custom domain to a domain in Amazon Route53

The problem: As of March 2020, Cognito User Pool domains created through CloudFormation do not return the ID of the CloudFront distribution backing them (see [here](https://github.com/aws-cloudformation/aws-cloudformation-coverage-roadmap/issues/356) and [here](https://github.com/aws-cloudformation/aws-cloudformation-coverage-roadmap/issues/58#issuecomment-539652016)). Because of this, you cannot link the domain to a custom domain you have in a Route53 hosted zone.
The solution: Deploy a Serverless Application Repository app which consists of a CloudFormation custom resource to do this for you!

![architecture](https://raw.githubusercontent.com/swoldemi/amazon-cognito-custom-domain-link/master/screenshots/architecture.png)

## Requirements
1. Your custom DNS name is hosted in Route 53
2. You have created an Amazon Certificate Manager certificate in us-east-1. This is the required region for CloudFront to be able to see the certificate
3. Your root domain (for instance, example.com), MUST have a valid A record in Route 53

## Usage
As stated above, there are two ways to use this Serverless Application Repo app. These examples assume you are creating a [`AWS::Cognito::UserPoolDomain`](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cognito-userpooldomain.html) in the same stack that you are using this custom resource. Note that, deleting and recreating a `AWS::Cognito::UserPoolDomain` can 15 minutes to fully create, 20 minutes to delete, and 1 HOUR for the deletion to fully propagate if you are frequently plan on attempting frequent creations and deletions of custom domains for your user pool's custom domain.

## Parameters
|Name           |Default           |Description                                                       |Required|                 
|---------------|------------------|------------------------------------------------------------------|--------|
|RegistryRegion |Function's Region |What AWS region should this Lambda function interact with ECR in? |False  -|


### 1. Serverless Application Model 
It is recommended that you use this custom resource as Nested Severless Application Repository app in your CloudFormation template.
(TODO)

### 2. CloudFormation
It is also possible to deploy this function seperately and link the Lambda function's ARN as the ServiceToken propery of the `AWS::CloudFormation::CustomResource`. Here is an example template (also available [here](./example-cfn-template.yaml))
```yaml
Parameters:
  Domain:
    Type: String
    Description: The domain name to use for both the User Pool and the Route 53 record (linked to CloudFront).
    Default: auth.simonw-aws.cloud # My example domain
  UserPoolId: # Remove if you are creating your User Pool in the same stack
    Type: String
    Description: ID of the Cognito User Pool that you are creating a custom domain for.
    Default: us-east-2_87lhunrJC # My example User Pool
  CertificateArn:
    Type: String
    Description: The ARN of the Amazon Certificate Manager certificate to be associate with your UserPoolDomain.
    Default: arn:aws:acm:us-east-1:273450712882:certificate/86e843a2-e0d3-496c-93b5-d9762da974f9 # My example certificate
  HostedZoneID:
    Type: String
    Description: The ID of the Route 53 hosted zone associated with your registered domain.
    Default: Z3I7ZCKH5Q5GZX # My hosted zone

Resources:
  UserPoolDomainRoute53Linker:
    Type: AWS::CloudFormation::CustomResource
    Version: 1.0.0
    Properties:
      # Lambda function in my AWS account
      ServiceToken: arn:aws:lambda:us-east-2:273450712882:function:amazon-cognito-custom-domain-link 
      HostedZoneID: !Ref HostedZoneID
      Domain: !Ref Domain

  UserPoolDomain:
    Type: AWS::Cognito::UserPoolDomain
    Properties:
      # Change this Ref if you are also creating your User Pool in the same Stack or using a Nested Stack
      UserPoolId: !Ref UserPoolId 
      Domain: !Ref Domain 
      CustomDomainConfig:
        CertificateArn: !Ref CertificateArn
```
To deploy this function from AWS GovCloud or regions in China, you must have an account with access to these regions. You must also deploy this function in the same region as your Cognito User Pool. This function is available in all regions that support AWS API Gateway, AWS Lambda, and Amazon Route53. If the table below is missing a region, please open a pull request!


|Region                                        |Click and Deploy                                                                                                                                 |
|----------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------|
|**US East (Ohio) (us-east-2)**                |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-east-2/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)     |
|**US East (N. Virginia) (us-east-1)**         |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-east-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)     |
|**US West (N. California) (us-west-1)**       |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-west-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)     |
|**US West (Oregon) (us-west-2)**              |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-west-2/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)     |
|**Asia Pacific (Hong Kong) (ap-east-1)**      |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-east-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)     |
|**Asia Pacific (Mumbai) (ap-south-1)**        |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-south-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)    |
|**Asia Pacific (Seoul) (ap-northeast-2)**     |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-northeast-2/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)|
|**Asia Pacific (Singapore)	(ap-southeast-1)** |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-southeast-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)|
|**Asia Pacific (Sydney) (ap-southeast-2)**    |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-southeast-2/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)|
|**Asia Pacific (Tokyo) (ap-northeast-1)**     |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-northeast-1?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link) |
|**Canada (Central)	(ca-central-1)**           |[![][sar-deploy]](https://deploy.serverlessrepo.app/ca-central-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)  |
|**EU (Frankfurt) (eu-central-1)**             |[![][sar-deploy]](https://deploy.serverlessrepo.app/eu-central-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)  |
|**EU (Ireland)	(eu-west-1)**                  |[![][sar-deploy]](https://deploy.serverlessrepo.app/eu-west-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)     |
|**EU (London) (eu-west-2)**                   |[![][sar-deploy]](https://deploy.serverlessrepo.app/eu-west-2/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)     |
|**EU (Paris) (eu-west-3)**                    |[![][sar-deploy]](https://deploy.serverlessrepo.app/eu-west-3/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)     |
|**EU (Stockholm) (eu-north-1)**               |[![][sar-deploy]](https://deploy.serverlessrepo.app/eu-north-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)    |
|**Middle East (Bahrain) (me-south-1)**        |[![][sar-deploy]](https://deploy.serverlessrepo.app/me-south-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)    |
|**South America (Sao Paulo) (sa-east-1)**     |[![][sar-deploy]](https://deploy.serverlessrepo.app/sa-east-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link)     |
|**AWS GovCloud (US-East) (us-gov-east-1)**    |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-gov-east-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link) |
|**AWS GovCloud (US-West) (us-gov-west-1)**    |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-gov-west-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-cognito-custom-domain-link) |

## Contributing
Have an idea for a feature to enhance this serverless application? Open an [issue](https://github.com/swoldemi/amazon-cognito-custom-domain-link/issues) or [pull request](https://github.com/swoldemi/amazon-cognito-custom-domain-link/pulls)!

### Development
This application has been developed, built, and tested against [Go 1.14](https://golang.org/dl/), the latest version of the [Serverless Application Model CLI](https://github.com/awslabs/aws-sam-cli), and the latest version of the [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html). A [Makefile](./Makefile) has been provided for convenience.

```
make check         # Run Go linting tools
make test          # Run Go tests
make build         # Build Go binary
make sam-package   # Package code and assets into S3 using SAM CLI
make sam-deploy    # Deploy application using SAM CLI
make sam-tail-logs # Tail the logs of the running Lambda function
make destroy       # Destroy the CloudFormation stack tied to the SAR app
```

### To Do
1. Support updating certificates on CloudFormation stack update by exposing a ACMCertificateARN parameter
2. Provide an option to block stack completion until DNS propagation is complete
3. Handle deletions and updates
  - DELETE: RecordSet - related to the CloudFront ALIAS 
  - UPDATE: Recordset - domain name changes related to the Cognito User Pool custom domain

## References
Using Your Own Domain for the Hosted UI - https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-add-custom-domain.html
## License
[Apache License 2.0](https://spdx.org/licenses/Apache-2.0.html)
