package main

import (
	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	log "github.com/sirupsen/logrus"
	"github.com/swoldemi/amazon-cognito-domain-distribution/pkg/function"
)

func main() {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("Error creating session: %v\n", err)
		return
	}

	cognitoSvc := cognitoidentityprovider.New(sess)
	lambda.Start(cfn.LambdaWrap(function.NewContainer(cognitoSvc).GetHandler()))
}
