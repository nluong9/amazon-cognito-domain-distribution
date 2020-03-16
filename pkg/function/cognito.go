package function

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	backoff "github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"
)

// ErrInvalidDomainState is a controlled error returned for the purpose
// of retrying with back off on the DescribeUserPoolDomain API.
var ErrInvalidDomainState = errors.New("domain does not have cloudfront distribution or has failed")

// GetPoolDistribution gets the DNS name of the CloudFront distribution associated
// with a Cognito User Pool custom domain. Because the domain does not have
// the DNS name of the CloudFront distribution immediately available, this method
// will retry, with back off, until it is available.
func (c *Container) GetPoolDistribution(ctx context.Context, domain string) (string, error) {
	var distribution string

	deleted, err := c.CheckPoolDeleted(ctx, domain)
	if err != nil {
		return "", err
	}
	if deleted {
		log.Infof("Pool %s has already been deleted", domain)
		return "", nil
	}

	operation := func() error {
		log.Infof("Describing user pool domain: %s", domain)
		input := &cognitoidentityprovider.DescribeUserPoolDomainInput{
			Domain: aws.String(domain),
		}
		if err := input.Validate(); err != nil {
			return err
		}
		output, err := c.CognitoIdentityProvider.DescribeUserPoolDomainWithContext(ctx, input)
		if err != nil {
			return err
		}
		distribution = aws.StringValue(output.DomainDescription.CloudFrontDistribution)
		status := aws.StringValue(output.DomainDescription.Status)
		desiredState := (status == cognitoidentityprovider.DomainStatusTypeCreating || status == cognitoidentityprovider.DomainStatusTypeActive)
		log.Debugf("Got user pool domain status: %s", status)

		if distribution == "" || !desiredState { // not ready OR reached DELETING or FAILED
			return ErrInvalidDomainState
		}
		return nil
	}

	if err := backoff.Retry(operation, backoff.NewExponentialBackOff()); err != nil {
		return "", err
	}
	log.Infof("Got CloudFront Distribution [%s]", distribution) // example: d111111abcdef8.cloudfront.net
	return distribution, nil
}

// CheckPoolDeleted checks if a User Pool domain has already been deleted or is in failure state.
func (c *Container) CheckPoolDeleted(ctx context.Context, domain string) (bool, error) {
	log.Infof("Verifying user pool has not been deleted yet: %s", domain)
	input := &cognitoidentityprovider.DescribeUserPoolDomainInput{
		Domain: aws.String(domain),
	}
	if err := input.Validate(); err != nil {
		return true, err
	}
	output, err := c.CognitoIdentityProvider.DescribeUserPoolDomainWithContext(ctx, input)
	if err != nil {
		return true, err
	}
	state := aws.StringValue(output.DomainDescription.Status)
	deletedState := (state == cognitoidentityprovider.DomainStatusTypeDeleting || state == cognitoidentityprovider.DomainStatusTypeFailed)
	if deletedState {
		return true, nil
	}
	return false, nil
}
