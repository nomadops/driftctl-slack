package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/logshell/driftctl-slack/driftctl"
	driftslack "github.com/logshell/driftctl-slack/slack"
)

func main() {
	// get the driftctl scan output filename from the environment
	driftctlJSON, err := os.LookupEnv("DRIFTCTL_JSON")
	if !err {
		log.Fatal("Environment variable DRIFTCTL_JSON does not exist")
	}

	// get the slack channel from the environment
	slackChannel, err := os.LookupEnv("SLACK_CHANNEL")
	if !err {
		log.Fatal("Environment variable SLACK_CHANNEL does not exist")
	}

	// Get the slack token from the environment
	slackToken, err := os.LookupEnv("SLACK_TOKEN")
	if !err {
		log.Fatal("Environment variable SLACK_TOKEN does not exist")
	}

	// Read the driftctl scan output file
	content, err1 := ioutil.ReadFile(driftctlJSON)
	if err1 != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Get Driftctl scan summary
	slackMessage, err1 := driftctl.ScanSummary(content)
	if err1 != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Send the driftctl scan summary output to slack
	driftslack.SendSummary(slackToken, slackChannel, slackMessage)
	if !err {
		log.Fatal("Error when running driftctl.ScanSummary", err)
	}
}
