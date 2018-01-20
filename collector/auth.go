package main

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var ErrMissingTwitterAuthTokens = errors.New("one or more authentication tokens missing")

type TwitterAuthDetails struct {
	ConsumerKey    string
	ConsumerSecret string
	Token          string
	Secret         string
}

func GetTwitterAuthDetails() (*TwitterAuthDetails, error) {

	// Fetch Twitter authentication tokens from AWS SSM Parameter Store

	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	// Set the AWS Region that the service clients should use
	cfg.Region = endpoints.EuWest1RegionID

	svc := ssm.New(cfg)
	req := svc.GetParametersRequest(&ssm.GetParametersInput{
		WithDecryption: aws.Bool(true),
		Names: []string{
			"/twitter/consumer-key",
			"/twitter/consumer-secret",
			"/twitter/token",
			"/twitter/secret",
		},
	})

	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	auth := &TwitterAuthDetails{}

	for _, param := range resp.Parameters {
		switch *param.Name {
		case "/twitter/consumer-key":
			auth.ConsumerKey = *param.Value
		case "/twitter/consumer-secret":
			auth.ConsumerSecret = *param.Value
		case "/twitter/token":
			auth.Token = *param.Value
		case "/twitter/secret":
			auth.Secret = *param.Value
		}
	}

	if len(auth.ConsumerSecret) < 1 || len(auth.ConsumerSecret) < 1 || len(auth.Token) < 1 || len(auth.Secret) < 1 {
		return nil, ErrMissingTwitterAuthTokens
	}

	return auth, nil

}
