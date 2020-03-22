package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/swoldemi/amazon-cognito-domain-distribution/pkg/function"
)

type mockCognitoClient struct {
	mock.Mock
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

// DescribeUserPoolDomainWithContext mocks the DescribeUserPoolDomain API.
func (_m *mockCognitoClient) DescribeUserPoolDomainWithContext(ctx aws.Context, input *cognitoidentityprovider.DescribeUserPoolDomainInput, opts ...request.Option) (*cognitoidentityprovider.DescribeUserPoolDomainOutput, error) {
	log.Debug("Mocking DescribeUserPoolDomainWithContext API")
	args := _m.Called(ctx, input)
	return args.Get(0).(*cognitoidentityprovider.DescribeUserPoolDomainOutput), args.Error(1)
}

const (
	defaultDistribution = "d111111abcdef8.cloudfront.net"
	domain              = "auth.example.com"
)

func TestHandler(t *testing.T) {
	createEvent := cfn.Event{
		RequestType:           cfn.RequestCreate,
		RequestID:             "c2664381-dabc-49d0-8347-bc2c01e12ea6",
		ResponseURL:           "http://pre-signed-S3-url-for-response",
		ResourceType:          "AWS::CloudFormation::CustomResource",
		PhysicalResourceID:    "",
		LogicalResourceID:     "UserPoolDomainDistribution",
		StackID:               "arn:aws:cloudformation:us-east-2:stack/amazon-cognito-domain-distribution/0921f5ae-0daf-4830-9a3c-ea1aa479a9d2",
		ResourceProperties:    map[string]interface{}{"Domain": domain},
		OldResourceProperties: map[string]interface{}{},
	}

	function.WithRetry = false
	tests := []struct {
		name         string
		event        cfn.Event
		wantErr      bool
		status       string
		distribution string
	}{

		{"DeletingState", createEvent, false, cognitoidentityprovider.DomainStatusTypeDeleting, ""},
		{"FailedState", createEvent, false, cognitoidentityprovider.DomainStatusTypeFailed, ""},
		{"ActiveState", createEvent, false, cognitoidentityprovider.DomainStatusTypeActive, defaultDistribution},
		{"UpdatingState", createEvent, false, cognitoidentityprovider.DomainStatusTypeUpdating, defaultDistribution},
		{"CreatingState", createEvent, false, cognitoidentityprovider.DomainStatusTypeCreating, defaultDistribution},
	}
	for _, tt := range tests {
		t.Logf("Running test: %+v", tt)
		t.Run(tt.name, func(t *testing.T) {
			cognitoSvc := new(mockCognitoClient)
			h := function.NewContainer(cognitoSvc).GetHandler()
			cognitoSvc.On("DescribeUserPoolDomainWithContext",
				context.Background(),
				&cognitoidentityprovider.DescribeUserPoolDomainInput{Domain: aws.String(domain)},
			).Return(
				&cognitoidentityprovider.DescribeUserPoolDomainOutput{
					DomainDescription: &cognitoidentityprovider.DomainDescriptionType{
						CloudFrontDistribution: aws.String(tt.distribution),
						Status:                 aws.String(tt.status),
					},
				}, nil)

			live := (tt.status == cognitoidentityprovider.DomainStatusTypeActive ||
				tt.status == cognitoidentityprovider.DomainStatusTypeUpdating ||
				tt.status == cognitoidentityprovider.DomainStatusTypeCreating)
			t.Logf("Live status is: %v", live)
			if live {
				cognitoSvc.On("DescribeUserPoolDomainWithContext",
					context.Background(),
					&cognitoidentityprovider.DescribeUserPoolDomainInput{Domain: aws.String(domain)},
				).Return(
					&cognitoidentityprovider.DescribeUserPoolDomainOutput{
						DomainDescription: &cognitoidentityprovider.DomainDescriptionType{
							CloudFrontDistribution: aws.String(tt.distribution),
							Status:                 aws.String(tt.status),
						},
					}, nil)
			}
			resource, outputs, err := h(context.Background(), tt.event)
			if tt.wantErr {
				require.NotNil(t, err)
			}
			if !tt.wantErr {
				require.Nil(t, err)
			}
			if live {
				require.Equal(t, defaultDistribution, outputs["CloudFrontDistributionDomainName"].(string))
			}
			if !live {
				require.Equal(t, "", outputs["CloudFrontDistributionDomainName"].(string))
			}
			t.Logf("Got physical resource ID: %s", resource)
			t.Logf("Got outputs: %v", outputs)
			cognitoSvc.AssertExpectations(t)
		})
	}
}
