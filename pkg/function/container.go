// Package function contains library units for the amazon-cognito-custom-domain-link Lambda function.
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
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrInvalidHostedZoneID ...
	ErrInvalidHostedZoneID = errors.New("invalid route 53 hosted zone ID provided")

	// ErrInvalidDomainName ...
	ErrInvalidDomainName = errors.New("invalid domain name provided")

	// ErrInvalidCreateRecordValue ...
	ErrInvalidCreateRecordValue = errors.New("invalid value provided for CreateRecord flag")
)

// Container contains the dependencies and business logic for the amazon-cognito-custom-domain-link Lambda function.
type Container struct {
	Route53                 route53iface.Route53API
	CognitoIdentityProvider cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

// NewContainer creates a new function Container.
func NewContainer(
	route53Svc route53iface.Route53API,
	cognitoSvc cognitoidentityprovideriface.CognitoIdentityProviderAPI,
) *Container {
	return &Container{
		Route53:                 route53Svc,
		CognitoIdentityProvider: cognitoSvc,
	}
}

var noop = make(map[string]interface{})

// GetHandler returns the function handler for the amazon-cognito-custom-domain-link.
// This custom resource expects two parameters to be set Route53HostedZoneID and CognitoUserPoolDomain.
func (c *Container) GetHandler() cfn.CustomResourceFunction {
	return func(ctx context.Context, event cfn.Event) (string, map[string]interface{}, error) {
		if event.PhysicalResourceID == "" {
			event.PhysicalResourceID = NewPhysicalResourceID(event)
		}
		log.Infof("Using physical resource ID: %s", event.PhysicalResourceID)

		domain, ok := event.ResourceProperties["Domain"].(string)
		if !ok {
			log.Errorf("Error during Domain lookup: %v", ErrInvalidDomainName)
			return event.PhysicalResourceID, noop, ErrInvalidDomainName
		}

		var hz string
		create, ok := event.ResourceProperties["CreateRecord"].(bool)
		if !ok {
			log.Errorf("Found non-boolean value for CreateRecord: %v", event.ResourceProperties)
			return event.PhysicalResourceID, noop, ErrInvalidCreateRecordValue
		}
		if create {
			hz, ok = event.ResourceProperties["HostedZoneID"].(string)
			if !ok {
				log.Errorf("Error during HostedZoneID lookup [CreateRecord is true]: %v", ErrInvalidHostedZoneID)
				return event.PhysicalResourceID, noop, ErrInvalidHostedZoneID
			}
		}

		distribution, err := c.GetPoolDistribution(ctx, domain)
		if err != nil {
			log.Errorf("Error during GetPoolDistributionID: %v", err)
			return event.PhysicalResourceID, noop, ErrInvalidDomainName
		}

		out := map[string]interface{}{
			"CloudFrontDistributionDomainName": distribution,
		}
		if !create {
			return event.PhysicalResourceID, out, nil
		}
		switch event.RequestType {
		case cfn.RequestCreate:
			return event.PhysicalResourceID, out, c.RunCreate(ctx, distribution, hz, domain)
		case cfn.RequestUpdate:
			return event.PhysicalResourceID, out, c.RunUpdate(ctx, distribution, hz, domain, event.OldResourceProperties)
		case cfn.RequestDelete:
			return event.PhysicalResourceID, out, c.RunDelete(ctx, distribution, hz, domain)
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
