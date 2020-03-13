package function

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/cfn"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrInvalidHostedZoneID ...
	ErrInvalidHostedZoneID = errors.New("resource: invalid route 53 hosted zone ID provided")

	// ErrInvalidDomainName ...
	ErrInvalidDomainName = errors.New("resource: invalid domain name provided")

	// ErrNotImplemented ...
	ErrNotImplemented = errors.New("resource: method not implemented")
)

// RunCreate runs the operations necessary for a CREATE event.
func (c *Container) RunCreate(ctx context.Context, event cfn.Event) error {
	hz, ok := event.ResourceProperties["HostedZoneID"].(string)
	if !ok {
		log.Errorf("Error during HostedZoneID lookup: %v", ErrInvalidHostedZoneID)
		return ErrInvalidHostedZoneID
	}
	domain, ok := event.ResourceProperties["Domain"].(string)
	if !ok {
		log.Errorf("Error during Domain lookup: %v", ErrInvalidDomainName)
		return ErrInvalidDomainName
	}

	distribution, err := c.GetPoolDistribution(ctx, domain)
	if err != nil {
		log.Errorf("Error during GetPoolDistributionID: %v", err)
		return err
	}
	if err := c.UpsertAlias(ctx, distribution, hz, domain); err != nil {
		log.Errorf("Error creating custom domain name")
		return err
	}
	return nil
}

// RunUpdate runs the operations necessary for a UPDATE event.
func (c *Container) RunUpdate(ctx context.Context, event cfn.Event) error {
	return ErrNotImplemented
}

// RunDelete runs the operations necessary for a DELETE event.
func (c *Container) RunDelete(ctx context.Context, event cfn.Event) error {
	return ErrNotImplemented
}
