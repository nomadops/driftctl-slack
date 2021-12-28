package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	c "github.com/nomadops/driftctl-slack/pkg/config"
	"github.com/nomadops/driftctl-slack/pkg/driftctl"
	drifts3 "github.com/nomadops/driftctl-slack/pkg/s3"
	driftslack "github.com/nomadops/driftctl-slack/pkg/slack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var err error
	// Apply zerolog global level, this will stop zerolog from logging any debug messages.
	debug := false
	// Apply log level in the beginning of the application
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Load the configuration from environment variables using Viper.
	env, err := c.LoadConfig()
	if err != nil {
		log.Fatal().
			Str("service", "driftctl-slack").
			Msg("Error loading viper config.")
		return
	}

	// Run driftctl scan command.
	message, err := driftctl.Run(env.StateBucket)
	if err != nil {
		log.Fatal().
			Str("service", "driftctl-slack").
			Msg("Error running driftctl scan.")
	}

	// Send driftctl scan summary output to slack
	err = driftslack.SendSummary(env.Token, env.Channel, message)
	if err != nil {
		log.Fatal().
			Str("service", "driftctl-slack").
			Msg("Error running driftctl.SendSummary.")
	}

	// Construct s3.PutObjectInput object
	msg, _ := json.Marshal(message)
	reader := strings.NewReader(string(msg))
	input := &s3.PutObjectInput{
		Bucket: aws.String(env.ScanBucket),
		Key:    aws.String(env.ScanReport),
		Body:   reader,
	}

	// Create a new S3 client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal().
			Str("service", "driftctl-slack").
			Msg("Error initializing S3 client.")
	}

	client := s3.NewFromConfig(cfg)

	// Copy the driftctl scan output file to S3 bucket.
	err = drifts3.PutFile(context.TODO(), client, input)
	if err != nil {
		log.Fatal().
			Str("service", "driftctl-slack").
			Msg("Error writing file to S3.")
	}
}
