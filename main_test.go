package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
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

func TestHandler(t *testing.T) {
	// TODO
	cognitoSvc := new(mockCognitoClient)
	h := function.NewContainer(cognitoSvc).GetHandler()
	_ = h
	_ = t
}
