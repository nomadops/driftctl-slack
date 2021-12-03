package driftslack

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/slack-go/slack"
)

// boldMarkup takes an array of strings and returns an array of strings with bold markdown
func boldMarkup(ary []string) []string {
	for i := 0; i < len(ary); i++ {
		ary[i] = "*" + ary[i] + "*"
	}
	return ary
}

// createSummaryBlockTable takes two arrays of strings to generate a summary table from driftctl scan output
func createSummaryBlockTable(statusLeft []string, statusRight []string) *slack.SectionBlock {
	// Format the keys
	statusLeftBold := boldMarkup(statusLeft)

	// Create a summary table from textblock objects
	leftText := strings.Join(statusLeftBold, "\n")
	rightText := strings.Join(statusRight, "\n")
	bodyLeft := slack.NewTextBlockObject("mrkdwn", leftText, false, false)
	bodyRight := slack.NewTextBlockObject("mrkdwn", rightText, false, false)
	fieldSlice := make([]*slack.TextBlockObject, 0)
	fieldSlice = append(fieldSlice, bodyLeft)
	fieldSlice = append(fieldSlice, bodyRight)
	fieldsSection := slack.NewSectionBlock(nil, fieldSlice, nil)
	return fieldsSection
}

// createSummaryMessage generates a slack message from driftctl scan output
func createSummaryMessage(summary map[string]int) slack.Message {
	var statusLeft []string
	var statusRight []string

	// iterate over driftctl scan summary
	for key, value := range summary {
		statusLeft = append(statusLeft, key)
		statusRight = append(statusRight, strconv.Itoa(value))
	}

	// Create a summary table from textblock objects
	fieldsSection := createSummaryBlockTable(statusLeft, statusRight)

	// Create a new header block
	headerText := slack.NewTextBlockObject("mrkdwn", "*Summary*", false, false)
	divSection := slack.NewDividerBlock()
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Create a slack message
	msg := slack.NewBlockMessage(
		headerSection,
		divSection,
		fieldsSection,
		divSection,
	)

	return msg
}

// SendSummary sends a slack message to a slack channel
func SendSummary(slackToken string, slackChannel string, slackMessage map[string]int) error {
	msg := createSummaryMessage(slackMessage)

	api := slack.New(slackToken)

	attachment := slack.Attachment{}

	attachment.Blocks = msg.Blocks
	channelID, timestamp, err := api.PostMessage(
		slackChannel,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionText(string("*driftctl scan report*"), false))
	if err != nil {
		log.Fatal("Error posting message to slack.", err)
		return err
	}

	fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)

	return nil
}
