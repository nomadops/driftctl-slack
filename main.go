package main

import (
	"io/ioutil"
	"os"

	"github.com/logshell/driftctl-slack/driftctl"
	driftslack "github.com/logshell/driftctl-slack/slack"
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
	driftctlJSON, err := os.LookupEnv("DRIFTCTL_JSON")
	if !err {
		log.Fatal().Msg("Environment variable DRIFTCTL_JSON does not exist")
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

	// Read driftctl scan output file
	content, err1 := ioutil.ReadFile(driftctlJSON)
	if err1 != nil {
		log.Fatal().
			Bool("Error while opening file:", err).
			Msg("")
	}

	// Get Driftctl scan summary
	slackMessage, err1 := driftctl.ScanSummary(content)
	if err1 != nil {
		log.Fatal().
			Bool("Error when opening file: ", err).
			Msg("")
	}
	log.Info().
		Str("service", "driftctl-slack").
		Int("total_resources", slackMessage["total_resources"]).
		Int("total_changed", slackMessage["total_changed"]).
		Int("total_unmanaged", slackMessage["total_unmanaged"]).
		Int("total_missing", slackMessage["total_missing"]).
		Int("total_managed", slackMessage["total_managed"]).
		Msg("Driftctl scan summary")

	// Send driftctl scan summary output to slack
	driftslack.SendSummary(slackToken, slackChannel, slackMessage)
	if !err {
		log.Fatal().
			Bool("Error when running driftctl.ScanSummary", err).
			Msg("")
	}
}

//
