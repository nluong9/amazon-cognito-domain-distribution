package function

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	log "github.com/sirupsen/logrus"
)

// See https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-aliastarget.html#cfn-route53-aliastarget-hostedzoneid
// This value will not change dynamically, and seems to have stayed the same since the release of CloudFront.
const cloudFrontHostedZone = "Z2FDTNDATAQYW2"

// UpsertAlias creates an ALIAS record for the current domain within Route53.
func (c *Container) UpsertAlias(ctx context.Context, distribution string, hz string, domain string) error {
	log.Infof("Creating alias of CloudFront distribution %s in hosted zone %s with domain %s", distribution, hz, domain)
	changes := []*route53.Change{
		{
			Action: aws.String(route53.ChangeActionUpsert),
			ResourceRecordSet: &route53.ResourceRecordSet{
				AliasTarget: &route53.AliasTarget{
					DNSName:              aws.String(distribution),
					EvaluateTargetHealth: aws.Bool(false),
					HostedZoneId:         aws.String(cloudFrontHostedZone),
				},
				Name: aws.String(domain),
				Type: aws.String(route53.RRTypeA),
			},
		},
	}
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: changes,
		},
		HostedZoneId: aws.String(hz),
	}
	if err := input.Validate(); err != nil {
		return err
	}
	if _, err := c.Route53.ChangeResourceRecordSetsWithContext(ctx, input); err != nil {
		return err
	}
	return nil
}

// DeleteAlias deletes an ALIAS record used to for the CloudFront
// distribution backing the Cognito User Pool customd omain
func (c *Container) DeleteAlias() error {
	return nil
}
