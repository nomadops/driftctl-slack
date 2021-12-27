package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/nomadops/driftctl-slack/pkg/driftctl"
	drifts3 "github.com/nomadops/driftctl-slack/pkg/s3"
	driftslack "github.com/nomadops/driftctl-slack/pkg/slack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Apply zerolog global level, this will stop zerolog from logging any debug messages.
	debug := false
	// Apply log level in the beginning of the application
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// get driftctl scan output filename from the environment
	scanReport, err := os.LookupEnv("SCAN_FILE")
	if !err {
		log.Fatal().Msg("Environment variable SCAN_FILE does not exist")
	}

	// get scan bucket name from the environment.
	scanBucket, err := os.LookupEnv("SCAN_BUCKET")
	if !err {
		log.Fatal().Msg("Environment variable SCAN_BUCKET does not exist")
	}

	// get state bucket name from the environment.
	stateBucket, err := os.LookupEnv("STATE_BUCKET")
	if !err {
		log.Fatal().Msg("Environment variable STATE_BUCKET does not exist")
	}

	// get slack channel from the environment
	channel, err := os.LookupEnv("CHANNEL")
	if !err {
		log.Fatal().Msg("Environment variable CHANNEL does not exist")
	}

	// Get slack token from the environment
	token, err := os.LookupEnv("TOKEN")
	if !err {
		log.Fatal().Msg("Environment variable TOKEN does not exist")
	}

	// Needs better variables
	message, err1 := driftctl.Run(stateBucket)
	if err1 != nil {
		log.Fatal().Msg("Error when running driftctl scan")
	}

	// Send driftctl scan summary output to slack
	driftslack.SendSummary(token, channel, message)
	if !err {
		log.Fatal().
			Bool("Error when running driftctl.ScanSummary", err).
			Msg("")
	}

	harry, _ := json.Marshal(message)
	reader := strings.NewReader(string(harry))
	input := &s3.PutObjectInput{
		Bucket: aws.String(scanBucket),
		Key:    aws.String(scanReport),
		Body:   reader,
	}

	// Create a new S3 client
	cfg, err1 := config.LoadDefaultConfig(context.TODO())
	if err1 != nil {
		log.Fatal().
			Bool("Error when running driftctl.ScanSummary", err).
			Msg("")
	}
	client := s3.NewFromConfig(cfg)

	// Copy the driftctl scan output file to S3 bucket.
	putOutput, err1 := drifts3.PutFile(context.TODO(), client, input)
	if err1 != nil {
		log.Fatal().
			Bool("Error when opening file: ", err).
			Msg("")
	}

	log.Info().
		Str("service", "driftctl-slack").
		Str("VersionId", *putOutput.VersionId).
		Msg("Report uploaded to S3")
}
