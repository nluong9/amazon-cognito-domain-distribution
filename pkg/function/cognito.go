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

// GetPoolDistribution kicks off the retrival of the distribution's DNS name.
func (c *Container) GetPoolDistribution(ctx context.Context, domain string, retry bool) (string, error) {
	var distribution string
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
		s := aws.StringValue(output.DomainDescription.Status)
		available := (s == cognitoidentityprovider.DomainStatusTypeCreating ||
			s == cognitoidentityprovider.DomainStatusTypeActive ||
			s == cognitoidentityprovider.DomainStatusTypeUpdating) // TODO: Does updating state consistently cause Updating state?

		log.Debugf("Got user pool domain status: %s", s)
		if distribution == "" || !available { // Not ready OR reached DELETING or FAILED
			return ErrInvalidDomainState
		}
		log.Infof("Got CloudFront Distribution [%s]", distribution) // Example: d111111abcdef8.cloudfront.net
		return nil
	}

	if retry {
		log.Info("Invoking GetPoolDistribution operation with retry")
		if err := backoff.Retry(operation, backoff.NewExponentialBackOff()); err != nil {
			return "", err
		}
		return distribution, nil
	}
	if err := operation(); err != nil {
		return "", err
	}
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
	s := aws.StringValue(output.DomainDescription.Status)
	deleted := (s == cognitoidentityprovider.DomainStatusTypeDeleting || s == cognitoidentityprovider.DomainStatusTypeFailed)
	if deleted {
		return true, nil
	}
	return false, nil
}
