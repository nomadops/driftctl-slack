package main

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/nomadops/driftctl-slack/driftctl"
	drifts3 "github.com/nomadops/driftctl-slack/s3"
	driftslack "github.com/nomadops/driftctl-slack/slack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Apply zerolog global level, this will stop zerolog from logging any debug messages
	debug := false
	// Apply log level in the beginning of the application
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// get driftctl scan output filename from the environment
	scanReport, err := os.LookupEnv("DRIFTCTL_JSON")
	if !err {
		log.Fatal().Msg("Environment variable DRIFTCTL_JSON does not exist")
	}

	// get scan bucket name from the environment
	scanBucket, err := os.LookupEnv("SCAN_BUCKET")
	if !err {
		log.Fatal().Msg("Environment variable SCAN_BUCKET does not exist")
	}

	// get state bucket name from the environment
	stateBucket, err := os.LookupEnv("STATE_BUCKET")
	if !err {
		log.Fatal().Msg("Environment variable STATE_BUCKET does not exist")
	}

	// get slack channel from the environment
	slackChannel, err := os.LookupEnv("SLACK_CHANNEL")
	if !err {
		log.Fatal().Msg("Environment variable SLACK_CHANNEL does not exist")
	}

	// Get slack token from the environment
	slackToken, err := os.LookupEnv("SLACK_TOKEN")
	if !err {
		log.Fatal().Msg("Environment variable SLACK_TOKEN does not exist")
	}

	// Needs better variables
	scanOutput, err1 := driftctl.Run(stateBucket, scanReport)
	if err1 != nil {
		log.Fatal().Msg("Error when running driftctl scan")
	}
	log.Info().Str("scan_output", string(scanOutput)).Msg("Scan output.")

	// // Read driftctl scan output file
	content, err1 := os.Open(scanReport)
	if err1 != nil {
		log.Fatal().
			Bool("Error while opening file:", err).
			Msg("")
	}
	defer content.Close()

	bob, _ := ioutil.ReadAll(content)
	// Get Driftctl scan summary
	slackMessage, err1 := driftctl.ScanSummary(bob)
	if err1 != nil {
		log.Fatal().
			Bool("Error when opening file: ", err).
			Msg("")
	}

	// Send driftctl scan summary output to slack
	driftslack.SendSummary(slackToken, slackChannel, slackMessage)
	if !err {
		log.Fatal().
			Bool("Error when running driftctl.ScanSummary", err).
			Msg("")
	}

	driftslack.Log(slackMessage)

	input := &s3.PutObjectInput{
		Bucket: aws.String(scanBucket),
		Key:    aws.String(scanReport),
		Body:   content,
	}

	// Create a new S3 client
	cfg, err1 := config.LoadDefaultConfig(context.TODO())
	if err1 != nil {
		log.Fatal().
			Bool("Error when running driftctl.ScanSummary", err).
			Msg("")
	}
	client := s3.NewFromConfig(cfg)

	george, err1 := drifts3.PutFile(context.TODO(), client, input)
	if err1 != nil {
		log.Fatal().
			Bool("Error when opening file: ", err).
			Msg("")
	}

	log.Info().Str("VersionId", *george.VersionId).Msg("Report uploaded to S3")
}
