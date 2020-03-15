![](https://codebuild.us-east-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiWnM4bGEzaDNFSXRUZFZPSmo4aDVVMjRBTDFOdVpRZ2kyM1pWSGdkV0ZaYjlCeTR2cjVrRmJWUXRuNFBpbFI0R1ZQU3pLbWQwUHJXT25tRVJUWW8zRjdRPSIsIml2UGFyYW1ldGVyU3BlYyI6IkozK01zeG5haDZQRVMvNTkiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)
[![][sar-logo]](https://serverlessrepo.aws.amazon.com/applications/arn:aws:serverlessrepo:us-east-1:273450712882:applications~amazon-cognito-custom-domain-link)

[sar-deploy]: https://img.shields.io/badge/Serverless%20Application%20Repository-Deploy%20Now-FF9900?logo=amazon%20aws&style=flat-square
[sar-logo]: https://img.shields.io/badge/Serverless%20Application%20Repository-View-FF9900?logo=amazon%20aws&style=flat-square

# Amazon Cognito Custom Domain Link
>A CloudFormation Custom Resource for automatically linking your Cognito User Pool's custom domain to a domain in Amazon Route 53

The problem: As of March 2020, Cognito User Pool domains created through CloudFormation do not return the DNS name of the CloudFront distribution backing them (see [here](https://github.com/aws-cloudformation/aws-cloudformation-coverage-roadmap/issues/356) and [here](https://github.com/aws-cloudformation/aws-cloudformation-coverage-roadmap/issues/58#issuecomment-539652016)). Because of this, you cannot link the domain to a custom domain you have in a Route53 hosted zone.

The solution: Deploy a Serverless Application Repository app which consists of a CloudFormation custom resource to do this for you! Similar to [`AWS::Route53::RecordSet`](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-recordset.html), this resource will create an alias target using a provided hosted zone ID and dns name. Alternatively, you can use this resource only for the purpose of returning the CloudFront distribution's DNS name. This provides for flexibility in the event that you would like to use [`AWS::Route53::RecordSetGroup`](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53-recordsetgroup.html)s in your template or configure additional attributes of the record set.

![architecture](https://raw.githubusercontent.com/swoldemi/amazon-cognito-custom-domain-link/master/screenshots/architecture.png)


## Requirements
1. Your custom DNS name is hosted in Route 53
2. You have created an Amazon Certificate Manager certificate in us-east-1. This is the required region for CloudFront to be able to see the certificate
3. Your zone apex (for instance, example.com or amazon.com) MUST have a valid A/AAAA record (IP address or aliased AWS resource) in Route 53
4. ***IF*** your domain is internationalized or uses emojis, you ***MUST*** [convert](https://www.punycoder.com/) it to [Punycode](https://en.wikipedia.org/wiki/Punycode) before passing it as a parameter to the template. This Lambda function will make no attempt at doing the conversion for you. See the following page for more details on how Route 53 handles internationalized/unicode domains ("Formatting Internationalized Domain Names" section): https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/DomainNameFormat.html

## Custom Resource Parameters
|Name           |Required |Description                                                                            |Default         |                 
|---------------|---------|---------------------------------------------------------------------------------------|-----------------
|Domain         |true     |The domain name to use for both the User Pool and the Route 53 record                  |auth.example.com|
|HostedZoneID   |false    |The ID of the Route 53 hosted zone associated with your registered domain.             |Z111111QQQQQQQ  |
|CreateRecord   |false    |(true or false) A flag to signify if this resource should also create the alias record |false            |

## Usage
These examples assume you are creating a [`AWS::Cognito::UserPoolDomain`](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cognito-userpooldomain.html) in the same stack that you are using this custom resource. Note that deleting and recreating a `AWS::Cognito::UserPoolDomain` can take 15 minutes to fully create, 20 minutes to delete, and 1 hour for the deletion to fully propagate through AWS if you are planning on attempting frequent creations and deletions. All examples are available in the [examples folder](./examples)

### 1. Serverless Application Model Template
It is recommended that you use this custom resource as a Severless Application Repository nested app in your CloudFormation template.
 - Example that will only return the DNS name of the CloudFront distribution: [here](./examples/sam/no-create-sam-template.yaml)
 - Example that also creates an alias record: [here](./examples/sam/sam-template.yaml)

### 2. CloudFormation Template
It is also possible to deploy this function seperately and link the Lambda function's ARN as the ServiceToken propery of the `AWS::CloudFormation::CustomResource`. The default values of the 3 parameters are sufficient to safely deploy the SAR application and retrieve the Lambda function's Amazon Resource Name. 
 - Example that will only return the DNS name of the CloudFront distribution: [here](./examples/cloudformation/no-create-cfn-template.yaml)
 - Example that also creates a domain: [here](./examples/cloudformation/cfn-template.yaml)

## Contributing
Have an idea for a feature to enhance this serverless application? Running into problems using it? Open an [issue](https://github.com/swoldemi/amazon-cognito-custom-domain-link/issues) or [pull request](https://github.com/swoldemi/amazon-cognito-custom-domain-link/pulls)!

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
Here are the Route 53 behaviors of this custom resource on each CloudFormation event. These only apply if the `CreateRecord` parameter is `true`. In all 3 event types, this resource will first get the DNS name of the CloudFormation distribution backing your Cognito User Pool domain and return it as an output attribute.
- `CREATE`
  1. Creates a new Route 53 Record Set (A Record) which will alias your custom domain against this distribution.
- `UPDATE`
  1. Compares the existing record with the new DNS name.
  2. No operation is performed if they are the same. 
  3. If they are the same, the existing record is deleted and a new record is created.  
- `DELETE`
  1. Check for existing record matching the provided parameters and deletes the existing Route 53 Record Set (A Record).
  2. Note: You create any traffic policies, this resource will not delete them. If you do not want to continue paying ($50/month) for the policies, remember to delete them.

### Return Values
#### Fn::GetAtt
The `Fn::GetAtt` intrinsic function returns a value for a specified attribute of this custom resource. This custom resource currently only returns one attribute.
For more information about using the Fn::GetAtt intrinsic function, see [Fn::GetAtt](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html).
- `CloudFrontDistributionDomainName`
    - The DNS name of the CloudFront distribution associated with the Cognito User Pool domain, such as d111111abcdef8.cloudfront.net.
  

### To Do
1. Provide an option to block stack `CREATE_COMPLETE` until DNS propagation is complete using the [route53.WaitUntilResourceRecordSetsChanged](https://pkg.go.dev/github.com/aws/aws-sdk-go/service/route53?tab=doc#Route53.WaitUntilResourceRecordSetsChanged) waiter.

## References
Using Your Own Domain for the Hosted UI - https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-add-custom-domain.html

What are some best practices for implementing AWS Lambda-backed custom resources with AWS CloudFormation? - https://aws.amazon.com/premiumsupport/knowledge-center/best-practices-custom-cf-lambda/

## License
[Apache License 2.0](https://spdx.org/licenses/Apache-2.0.html)
