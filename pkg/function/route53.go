package function

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	log "github.com/sirupsen/logrus"
)

// CloudFrontHostedZone is the hosted zone used for CloudFront domains.
// See https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-aliastarget.html#cfn-route53-aliastarget-hostedzoneid
// This value will not change dynamically and has alwas had this value.
const CloudFrontHostedZone = "Z2FDTNDATAQYW2"

// Domain is the domain name and Route53 hosted zone of a domain record
type Domain struct {
	HostedZoneID string
	Name         string
}

// CreateAlias creates an ALIAS record for the current domain within Route53.
func (c *Container) CreateAlias(ctx context.Context, distribution string, domain *Domain) error {
	log.Infof("Creating alias of CloudFront distribution %s in hosted zone %s with domain %s", distribution, domain.HostedZoneID, domain.Name)
	input, err := NewRecordSetChange(distribution, domain, route53.ChangeActionCreate)
	if err != nil {
		return err
	}
	if _, err := c.Route53.ChangeResourceRecordSetsWithContext(ctx, input); err != nil {
		return err
	}
	return nil
}

// UpsertAlias updates an ALIAS record used to for the CloudFront
// distribution backing the Cognito User Pool custom domain.
func (c *Container) UpsertAlias(ctx context.Context, distribution string, od, nd *Domain) error {
	if od.Name == nd.Name && od.HostedZoneID == nd.HostedZoneID { // old domain and new domain are the same
		log.Warnf("Existing record [%s:%s] and new record [%s:%s] are the same, no change", od.HostedZoneID, od.Name, nd.HostedZoneID, nd.Name)
		return nil
	}

	log.Infof("Updating alias of CloudFront distribution %s in from [%s:%s] to [%s:%s]", distribution, od.HostedZoneID, od.Name, nd.HostedZoneID, nd.Name)

	// Check for the existing CloudFront Alias in the current distribution, associated with the current User Pool domain
	record, err := c.FindCloudFrontDistributionAliasRecordSet(ctx, distribution, od)
	if err != nil {
		return err
	}
	if record == nil {
		log.Warnf("Route 53 record set in hosted zone %s with name %s was not found; will not attempt deletetion", od.HostedZoneID, od.Name)
	}
	if record != nil {
		// Delete the existing record
		ed := &Domain{
			HostedZoneID: aws.StringValue(record.AliasTarget.HostedZoneId),
			Name:         aws.StringValue(record.Name),
		}
		if err := c.DeleteAlias(ctx, distribution, ed); err != nil {
			return err
		}
	}

	// Update/Create the new record
	input, err := NewRecordSetChange(distribution, nd, route53.ChangeActionUpsert)
	if err != nil {
		return err
	}
	if _, err := c.Route53.ChangeResourceRecordSetsWithContext(ctx, input); err != nil {
		return err
	}
	return nil
}

// DeleteAlias deletes an ALIAS record used to for the CloudFront
// distribution backing the Cognito User Pool custom domain.
func (c *Container) DeleteAlias(ctx context.Context, distribution string, domain *Domain) error {
	log.Infof("Deleting alias of CloudFront distribution %s in hosted zone %s with domain %s", distribution, domain.HostedZoneID, domain.Name)
	input, err := NewRecordSetChange(distribution, domain, route53.ChangeActionDelete)
	if err != nil {
		return err
	}
	if _, err := c.Route53.ChangeResourceRecordSetsWithContext(ctx, input); err != nil {
		return err
	}
	return nil
}

// NewRecordSetChange is a helper for creating the ChangeResourceRecordSetsInput.
func NewRecordSetChange(distribution string, domain *Domain, action string) (*route53.ChangeResourceRecordSetsInput, error) {
	changes := []*route53.Change{
		{
			Action: aws.String(action),
			ResourceRecordSet: &route53.ResourceRecordSet{
				AliasTarget: &route53.AliasTarget{
					DNSName:              aws.String(distribution),
					EvaluateTargetHealth: aws.Bool(false),
					HostedZoneId:         aws.String(CloudFrontHostedZone),
				},
				Name: aws.String(domain.Name),
				Type: aws.String(route53.RRTypeA),
			},
		},
	}
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: changes,
		},
		HostedZoneId: aws.String(domain.HostedZoneID),
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return input, nil
}

// FindCloudFrontDistributionAliasRecordSet finds if there is an A or AAAA record set containing
// the existing CloudFront distribution alias resource.
// `distribution` should be the DNS name of a CloudFront Distribution (e.e. randomletters.cloudfront.net)
func (c *Container) FindCloudFrontDistributionAliasRecordSet(ctx context.Context, distribution string, domain *Domain) (*route53.ResourceRecordSet, error) {
	var (
		exists bool
		record *route53.ResourceRecordSet
	)
	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(domain.HostedZoneID),
		StartRecordName: aws.String(checkdot(domain.Name)),
		StartRecordType: aws.String(route53.RRTypeA),
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	distribution = checkdot(distribution)
	pager := func(out *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
		for _, rr := range out.ResourceRecordSets {
			sameDistribution := checkdot(aws.StringValue(rr.AliasTarget.DNSName)) == distribution
			if aws.StringValue(rr.AliasTarget.HostedZoneId) == CloudFrontHostedZone && sameDistribution {
				exists = true
				record = rr
				return true // Found what we're looking for, exit
			}
		}
		return lastPage
	}
	if err := c.Route53.ListResourceRecordSetsPagesWithContext(ctx, input, pager); err != nil {
		return nil, err
	}
	if !exists {
		log.Warnf("Couldn't find A record matching CloudFront distribution, attempting search for AAAA records")
		input.StartRecordType = aws.String(route53.RRTypeAaaa)
		if err := c.Route53.ListResourceRecordSetsPagesWithContext(ctx, input, pager); err != nil {
			return nil, err
		}
	}
	return record, nil
}

// check if a given DNS name ends with a trailing dot
func checkdot(dn string) string {
	if strings.HasSuffix(dn, ".") {
		return dn
	}
	return fmt.Sprintf("%s.", dn)
}
