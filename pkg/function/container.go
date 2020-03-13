// Package function contains library units for the amazon-cognito-custom-domain-link Lambda function.
package function

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	log "github.com/sirupsen/logrus"
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
			event.PhysicalResourceID = generatePRID(event)
		}
		log.Infof("Using physical resource ID: %s", event.PhysicalResourceID)

		switch event.RequestType {
		case cfn.RequestCreate:
			return event.PhysicalResourceID, noop, c.RunCreate(ctx, event)
		case cfn.RequestUpdate, cfn.RequestDelete:
			log.Info("Request type is not CREATE; no operation")
			return event.PhysicalResourceID, noop, nil
		}
		return event.PhysicalResourceID, noop, fmt.Errorf("got invalid requet type: %s", event.RequestType)
	}
}

// Credit to @dweomer
// https://github.com/dweomer/aws-cloudformation-keypair/blob/master/aws/ec2/keypair/resource.go#L131-L145
func generatePRID(event cfn.Event) string {
	rns := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	gen := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd := make([]byte, 12)
	for i := range rnd {
		rnd[i] = rns[gen.Intn(len(rns))]
	}
	stack := strings.Split(event.StackID, "/")[1]
	return fmt.Sprintf("%s-%s-%s", stack, event.LogicalResourceID, rnd)
}
