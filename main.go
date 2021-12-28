package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/nomadops/driftctl-slack/pkg/driftctl"
	drifts3 "github.com/nomadops/driftctl-slack/pkg/s3"
	driftslack "github.com/nomadops/driftctl-slack/pkg/slack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config is a struct that contains the configuration environment for the application.
type Config struct {
	scanReport  string `mapstructure:"SCAN_FILE"`
	scanBucket  string `mapstructure:"SCAN_BUCKET"`
	stateBucket string `mapstructure:"STATE_BUCKET"`
	channel     string `mapstructure:"CHANNEL"`
	token       string `mapstructure:"TOKEN"`
}

// LoadConfig loads the configuration from the environment.
func LoadConfig() (config Config, err error) {
	viper.SetConfigType("env")
	viper.BindEnv("SCAN_FILE")
	viper.BindEnv("SCAN_BUCKET")
	viper.BindEnv("STATE_BUCKET")
	viper.BindEnv("CHANNEL")
	viper.BindEnv("TOKEN")
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func main() {
	// Apply zerolog global level, this will stop zerolog from logging any debug messages.
	debug := false
	// Apply log level in the beginning of the application
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// get driftctl scan output filename from the environment
	env, err1 := LoadConfig()
	if err1 != nil {
		log.Fatal().
			Str("service", "driftctl-slack").
			Msg("Error when loading viper config: ")
		return
	}

	// Needs better variables
	message, err1 := driftctl.Run(env.stateBucket)
	if err1 != nil {
		log.Fatal().Msg("Error when running driftctl scan")
	}

	// Send driftctl scan summary output to slack
	err := driftslack.SendSummary(env.token, env.channel, message)
	if err != nil {
		log.Fatal().
			Str("service", "driftctl-slack").
			Msg("Error when running driftctl.ScanSummary")
	}

	harry, _ := json.Marshal(message)
	reader := strings.NewReader(string(harry))
	input := &s3.PutObjectInput{
		Bucket: aws.String(env.scanBucket),
		Key:    aws.String(env.scanReport),
		Body:   reader,
	}

	// Create a new S3 client
	cfg, err1 := config.LoadDefaultConfig(context.TODO())
	if err1 != nil {
		log.Fatal().
			Str("service", "driftctl-s3").
			Msg("Error when running driftctl.ScanSummary")
	}

	client := s3.NewFromConfig(cfg)

	// Copy the driftctl scan output file to S3 bucket.
	err2 := drifts3.PutFile(context.TODO(), client, input)
	if err2 != nil {
		log.Fatal().
			Str("service", "driftctl-s3").
			Msg("Error writing file to S3:")
	}
}
