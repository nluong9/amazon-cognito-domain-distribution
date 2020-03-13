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
// of retrying with back off on the DescribeUserPoolDomain API
var ErrInvalidDomainState = errors.New("domain does not have cloudfront distribution or has failed")

// GetPoolDistribution gets CloudFront Distribution associated
// with a Cognito User Pool Custom Domain. Because the domain does not have
// the DNS name of the CloudFront distribution immediately available, this method
// will retry, with back off, until it is available.
func (c *Container) GetPoolDistribution(ctx context.Context, domain string) (string, error) {
	var distribution string

	operation := func() error {
		log.Infof("Describing user pool domain: %s", domain)
		input := &cognitoidentityprovider.DescribeUserPoolDomainInput{
			Domain: aws.String(domain),
		}
		if err := input.Validate(); err != nil {
			return err
		}
		output, err := c.CognitoIdentityProvider.DescribeUserPoolDomain(input)
		if err != nil {
			return err
		}
		distribution = aws.StringValue(output.DomainDescription.CloudFrontDistribution)
		status := aws.StringValue(output.DomainDescription.Status)
		isValidStatus := (status == cognitoidentityprovider.DomainStatusTypeCreating || status == cognitoidentityprovider.DomainStatusTypeActive)
		log.Debugf("Got user pool domain status: %s", status)

		if distribution == "" || !isValidStatus {
			return ErrInvalidDomainState
		}
		return nil
	}

	if err := backoff.Retry(operation, backoff.NewExponentialBackOff()); err != nil {
		return "", err
	}
	log.Infof("Got CloudFront Distribution [%s]", distribution)
	return distribution, nil
}
