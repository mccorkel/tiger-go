package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadDefaultConfig() (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO())
}
