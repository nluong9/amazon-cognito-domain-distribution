![](https://codebuild.us-east-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiRXh5VkFmNmdIeUtxbFNzVHBML0pLck1zZWxYeERoSTZybzVabXBSOWlpWTFPS0Z0bG1POXY5RGYvUUNvQTAwNmhIUWF1NkJxL2JuOHlsN0IvUzdNejNVPSIsIml2UGFyYW1ldGVyU3BlYyI6ImVidEJscmVZeHRZTm12L08iLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)
[![][sar-logo]](https://serverlessrepo.aws.amazon.com/applications/arn:aws:serverlessrepo:us-east-1:273450712882:applications~amazon-cognito-domain-distribution)

[sar-deploy]: https://img.shields.io/badge/Serverless%20Application%20Repository-Deploy%20Now-FF9900?logo=amazon%20aws&style=flat-square
[sar-logo]: https://img.shields.io/badge/Serverless%20Application%20Repository-View-FF9900?logo=amazon%20aws&style=flat-square

# Amazon Cognito Custom Domain Link
>An AWS CloudFormation Custom Resource for retrieving the DNS name of the Amazon CloudFront distribution backing your Amazon Cognito User Pool's custom domain

The problem: As of March 2020, Cognito User Pool domains created through CloudFormation do not return the DNS name of the CloudFront distribution backing them (see [here](https://github.com/aws-cloudformation/aws-cloudformation-coverage-roadmap/issues/356) and [here](https://github.com/aws-cloudformation/aws-cloudformation-coverage-roadmap/issues/58#issuecomment-539652016)). Because of this, you cannot link the domain to a custom domain you have in a Route53 hosted zone via CloudFormation.

The solution: Deploy a Serverless Application Repository app which consists of a CloudFormation custom resource to help you do this! This resource will return the CloudFront distribution's DNS name and you can use it in a [`AWS::Route53::RecordSet`](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-recordset.html) delcaration. See [here](./example-sam-template.yaml) for an example.

![architecture](https://raw.githubusercontent.com/swoldemi/amazon-cognito-domain-distribution/master/screenshots/architecture.png)

## Requirements
1. Your custom domain name is hosted in Route 53
2. You have created an Amazon Certificate Manager certificate in us-east-1. This is the required region for CloudFront to be able to see the certificate
3. Your zone apex (for instance, example.com or amazon.com) MUST have a valid A/AAAA record (IP address or aliased AWS resource) in Route 53
4. ***IF*** your domain is internationalized or uses emojis, you ***MUST*** [convert](https://www.punycoder.com/) it to [Punycode](https://en.wikipedia.org/wiki/Punycode) before passing it as a parameter to the template. This Lambda function will make no attempt at doing the conversion for you. See the following page for more details on how Route 53 handles internationalized/unicode domains ("Formatting Internationalized Domain Names" section): https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/DomainNameFormat.html

## Custom Resource Parameters
|Name           |Required |Description                                                           |                 
|---------------|---------|----------------------------------------------------------------------|
|Domain         |true     |The domain name to use for both the User Pool and the Route 53 record |

This should not be mistaken for SAR application parameters. The SAR application takes no parameters, but the custom resource (Lambda function) that the SAR application deploys takes 1 parameter.

## Custom Resource Return Values
#### Fn::GetAtt
The `Fn::GetAtt` intrinsic function returns a value for a specified attribute of this custom resource. This custom resource currently only returns one attribute.
For more information about using the Fn::GetAtt intrinsic function, see [Fn::GetAtt](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html).
- `CloudFrontDistributionDomainName`
    - The DNS name of the CloudFront distribution associated with the Cognito User Pool domain, such as d111111abcdef8.cloudfront.net.

![output-example](https://raw.githubusercontent.com/swoldemi/amazon-cognito-domain-distribution/master/screenshots/output.png)


## Usage

###  Serverless Application Model Template
It is recommended that you use this custom resource as a Severless Application Repository nested app.

The provided example assumes you are creating a [`AWS::Cognito::UserPoolDomain`](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cognito-userpooldomain.html) in the same stack that you are using this custom resource. Note that deleting and recreating a `AWS::Cognito::UserPoolDomain` can take 15 minutes to fully create, 20 minutes to delete, and 1 hour for the deletion to fully propagate through AWS if you are planning on attempting frequent creations and deletions. It is also possible to deploy the Lambda function template seperately and interact with vanillia CloudFormation instead.

The example template is available [here](./example-sam-template.yaml).

## Contributing
Have an idea for a feature to enhance this serverless application? Running into problems using it? Open an [issue](https://github.com/swoldemi/amazon-cognito-domain-distribution/issues) or [pull request](https://github.com/swoldemi/amazon-cognito-domain-distribution/pulls)!

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

### Behvaior
- Stack `CREATE`
  1. Returns the DNS name of the CloudFormation distribution.
- Stack `UPDATE`
  1. Returns the DNS name of the CloudFormation distribution. Returns an empty string ("") if the UPDATE involved deleting the User Pool domain.
- Stack `DELETE`
  1. No operation is performed.

### To Do
1. Expose configuration for cross-region Cognito interactions. As of March 2020, Amazon Cognito is available in 12 commercial regions. 

## References
Using Your Own Domain for the Hosted UI - https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-add-custom-domain.html

What are some best practices for implementing AWS Lambda-backed custom resources with AWS CloudFormation? - https://aws.amazon.com/premiumsupport/knowledge-center/best-practices-custom-cf-lambda/

## License
[Apache License 2.0](https://spdx.org/licenses/Apache-2.0.html)
