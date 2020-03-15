package function

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrInvalidExistingDomain ...
	ErrInvalidExistingDomain = errors.New("unable to find previous domain")

	// ErrInvalidExistingHostedZone ...
	ErrInvalidExistingHostedZone = errors.New("unable to find previous hosted zone")
)

// RunCreate runs the operations necessary for a CREATE event.
func (c *Container) RunCreate(ctx context.Context, distribution string, hz, domain string) error {
	if err := c.CreateAlias(ctx, distribution, &Domain{HostedZoneID: hz, Name: domain}); err != nil {
		log.Errorf("Error creating custom domain name")
		return err
	}
	return nil
}

// RunUpdate runs the operations necessary for a UPDATE event.
// FIXME: RunUpdate currently does not account for DNS propagation. If a resource is created
// and updated quickly, the state of the update may not match desired behavior.
// The idempotence of Route 53 may be enough to account for this.
func (c *Container) RunUpdate(ctx context.Context, distribution, hz, domain string, oldProperities map[string]interface{}) error {
	oldDomain, ok := oldProperities["Domain"].(string)
	if !ok {
		return ErrInvalidExistingDomain
	}
	oldHostedZoneID, ok := oldProperities["HostedZoneID"].(string)
	if !ok {
		return ErrInvalidExistingDomain
	}
	od := &Domain{
		HostedZoneID: oldHostedZoneID,
		Name:         oldDomain,
	}
	nd := &Domain{
		HostedZoneID: hz,
		Name:         domain,
	}
	if err := c.UpsertAlias(ctx, distribution, od, nd); err != nil {
		log.Errorf("Error updating custom domain name")
		return err
	}
	return nil
}

// RunDelete runs the operations necessary for a DELETE event.
func (c *Container) RunDelete(ctx context.Context, distribution, hz, domain string) error {
	if err := c.DeleteAlias(ctx, distribution, &Domain{HostedZoneID: hz, Name: domain}); err != nil {
		log.Errorf("Error deleting custom domain name")
		return err
	}
	return nil
}
