package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/swoldemi/amazon-cognito-custom-domain-link/pkg/function"
)

type mockRoute53Client struct {
	mock.Mock
	route53iface.Route53API
}

type mockCognitoClient struct {
	mock.Mock
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

// ChangeResourceRecordSetsWithContext mocks the ChangeResourceRecordSets API.
func (_m *mockRoute53Client) ChangeResourceRecordSetsWithContext(ctx aws.Context, input *route53.ChangeResourceRecordSetsInput, opts ...request.Option) (*route53.ChangeResourceRecordSetsOutput, error) {
	log.Debug("Mocking ChangeResourceRecordSetsWithContext API")
	args := _m.Called(ctx, input)
	return args.Get(0).(*route53.ChangeResourceRecordSetsOutput), args.Error(1)
}

// DescribeRepositoriesPagesWithContext mocks the DescribeRepositoriesPagesWithContext ECR API endpoint.
func (_m *mockRoute53Client) DescribeRepositoriesPagesWithContext(ctx aws.Context, input *route53.ListResourceRecordSetsInput, fn func(*route53.ListResourceRecordSetsOutput, bool) bool, opts ...request.Option) error {
	log.Debug("Mocking DescribeRepositoriesPagesWithContext API")
	args := _m.Called(ctx, input, fn, opts)
	return args.Error(0)
}

// DescribeUserPoolDomainWithContext mocks the DescribeUserPoolDomain API.
func (_m *mockCognitoClient) DescribeUserPoolDomainWithContext(ctx aws.Context, input *cognitoidentityprovider.DescribeUserPoolDomainInput, opts ...request.Option) (*cognitoidentityprovider.DescribeUserPoolDomainOutput, error) {
	log.Debug("Mocking DescribeUserPoolDomainWithContext API")
	args := _m.Called(ctx, input)
	return args.Get(0).(*cognitoidentityprovider.DescribeUserPoolDomainOutput), args.Error(1)
}

func TestHandler(t *testing.T) {
	// TODO
	route53Svc := new(mockRoute53Client)
	cognitoSvc := new(mockCognitoClient)
	h := function.NewContainer(route53Svc, cognitoSvc).GetHandler()
	_ = h
	_ = t
}
