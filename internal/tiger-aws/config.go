package tigeraws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func NewConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"), // Default region for the project
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return cfg
}
