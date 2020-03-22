// Package function contains library units for the amazon-cognito-domain-distribution Lambda function.
package function

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	log "github.com/sirupsen/logrus"
)

// WithRetry is a hacky global flag to make it possible test GetPoolDistribution without retry. FIXME.
var WithRetry = true

// ErrInvalidDomainName ...
var ErrInvalidDomainName = errors.New("invalid domain name provided")

// Container contains the dependencies and business logic for the amazon-cognito-domain-distribution Lambda function.
type Container struct {
	CognitoIdentityProvider cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

// NewContainer creates a new function Container.
func NewContainer(cognitoSvc cognitoidentityprovideriface.CognitoIdentityProviderAPI) *Container {
	return &Container{
		CognitoIdentityProvider: cognitoSvc,
	}
}

// GetHandler returns the function handler for the amazon-cognito-domain-distribution.
// This custom resource expects two parameters to be set Route53HostedZoneID and CognitoUserPoolDomain.
func (c *Container) GetHandler() cfn.CustomResourceFunction {
	return func(ctx context.Context, event cfn.Event) (string, map[string]interface{}, error) {
		out := map[string]interface{}{
			"CloudFrontDistributionDomainName": "",
		}
		log.Infof("Got resource properties: %v", event.ResourceProperties)
		if event.PhysicalResourceID == "" {
			event.PhysicalResourceID = NewPhysicalResourceID(event)
		}
		log.Infof("Using physical resource ID: %s", event.PhysicalResourceID)

		domain, ok := event.ResourceProperties["Domain"].(string)
		if !ok {
			log.Errorf("Error during Domain lookup: %v", ErrInvalidDomainName)
			return event.PhysicalResourceID, out, ErrInvalidDomainName
		}
		log.Infof("Got Cognito user pool domain name: %s", domain)

		switch event.RequestType {
		case cfn.RequestCreate, cfn.RequestUpdate:
			deleted, err := c.CheckPoolDeleted(ctx, domain)
			if err != nil {
				log.Errorf("Error during CheckPoolDeleted: %v", err)
				return event.PhysicalResourceID, out, err
			}
			if deleted {
				log.Infof("Pool %s has already been deleted", domain)
				return event.PhysicalResourceID, out, nil
			}
			distribution, err := c.GetPoolDistribution(ctx, domain, WithRetry)
			if err != nil {
				log.Errorf("Error during GetPoolDistribution: %v", err)
				return event.PhysicalResourceID, out, err
			}
			out["CloudFrontDistributionDomainName"] = distribution
			return event.PhysicalResourceID, out, nil
		case cfn.RequestDelete:
			log.Infof("Got DELETE event; no supported operation to perform")
			return event.PhysicalResourceID, out, nil
		}
		return event.PhysicalResourceID, out, fmt.Errorf("got invalid request type: %s", event.RequestType)
	}
}

// NewPhysicalResourceID creates a new physical resource ID.
// Credit to @dweomer
// https://github.com/dweomer/aws-cloudformation-keypair/blob/master/aws/ec2/keypair/resource.go#L131-L145
func NewPhysicalResourceID(event cfn.Event) string {
	rns := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	gen := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd := make([]byte, 12)
	for i := range rnd {
		rnd[i] = rns[gen.Intn(len(rns))]
	}
	stack := strings.Split(event.StackID, "/")[1]
	return fmt.Sprintf("%s-%s-%s", stack, event.LogicalResourceID, rnd)
}
