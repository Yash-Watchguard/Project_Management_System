package config

import (
	"context"

	lambdaconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func GetDyanoDbCliebt() (*dynamodb.Client) {
	cfg, err := lambdaconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("error in load config")
	}

	return dynamodb.NewFromConfig(cfg)
}
