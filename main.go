package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/slack-go/slack"
)

// DriftCtlOutput struct is the json output of `driftctl scan`
type DriftCtlOutput struct {
	Alerts struct {
		_ []struct {
			Message string `json:"message"`
		} `json:""`
	} `json:"alerts"`
	Coverage    int64       `json:"coverage"`
	Differences interface{} `json:"differences"`
	Managed     []struct {
		ID     string `json:"id"`
		Source struct {
			InternalName string `json:"internal_name"`
			Namespace    string `json:"namespace"`
			Source       string `json:"source"`
		} `json:"source"`
		Type string `json:"type"`
	} `json:"managed"`
	Missing []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"missing"`
	ProviderName    string `json:"provider_name"`
	ProviderVersion string `json:"provider_version"`
	Summary         struct {
		TotalChanged   int64 `json:"total_changed"`
		TotalManaged   int64 `json:"total_managed"`
		TotalMissing   int64 `json:"total_missing"`
		TotalResources int64 `json:"total_resources"`
		TotalUnmanaged int64 `json:"total_unmanaged"`
	} `json:"summary"`
	Unmanaged []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"unmanaged"`
}

func main() {
	content, err := ioutil.ReadFile(os.Getenv("DRIFTCTL_JSON"))
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var driftctloutput DriftCtlOutput
	err = json.Unmarshal(content, &driftctloutput)
	if err != nil {
		log.Fatal("Error when unmarshalling: ", err)
	}

	summary := driftctloutput.Summary

	bob := reflect.ValueOf(summary)
	typeOfS := bob.Type()

	var statusLeft []string
	var statusRight []string

	for i := 0; i < bob.NumField(); i++ {
		statusLeft = append(statusLeft, typeOfS.Field(i).Name)
		var george int64 = bob.Field(i).Interface().(int64)
		statusRight = append(statusRight, strconv.Itoa(int(george)))
	}

	for i := 0; i < len(statusLeft); i++ {
	}

	leftText := strings.Join(statusLeft, "\n")
	rightText := strings.Join(statusRight, "\n")

	headerText := slack.NewTextBlockObject("mrkdwn", "driftctl report", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	bodyLeft := slack.NewTextBlockObject("mrkdwn", leftText, false, false)
	bodyRight := slack.NewTextBlockObject("mrkdwn", rightText, false, false)

	fieldSlice := make([]*slack.TextBlockObject, 0)
	fieldSlice = append(fieldSlice, bodyLeft)
	fieldSlice = append(fieldSlice, bodyRight)

	fieldsSection := slack.NewSectionBlock(nil, fieldSlice, nil)

	msg := slack.NewBlockMessage(
		headerSection,
		fieldsSection,
	)

	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))

	api := slack.New(os.Getenv("SLACK_TOKEN"))

	attachment := slack.Attachment{}

	attachment.Blocks = msg.Blocks
	channelID, timestamp, err := api.PostMessage(
		"#gitops",
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionText(string("driftctl report"), false))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
