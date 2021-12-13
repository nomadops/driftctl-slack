#!/bin/bash

# Custom JSON logging function
log() {
  # Stolen from https://stegard.net/2021/07/how-to-make-a-shell-script-log-json-messages/
  jq --arg message "${1}" --arg hostname "${HOSTNAME}" --arg service "${CMD}[${PID}]" -nc '{"@timestamp": now|strftime("%Y-%m-%dT%H:%M:%S%z"),"hostname": $hostname, "service": $service, "message": $message}'
}

# Exit if the BUCKET_NAME environment variable is not set
if [[ ${STATE_BUCKET_NAME:-"unset"} == "unset" ]]; then
  log "Environment variable STATE_BUCKET_NAME is not set"
  exit 1
fi

# Exit if the SCAN_BUCKET_NAME environment variable is not set
if [[ ${SCAN_BUCKET_NAME:-"unset"} == "unset" ]]; then
  log "Environment variable SCAN_BUCKET_NAME is not set"
  exit 1
fi

# Exit if the DRIFTCTL_JSON environment variable is not set
if [[ ${DRIFTCTL_JSON:-"unset"} == "unset" ]]; then
  log "Environment variable DRIFTCTL_JSON is not set"
  exit 1
fi

# Exit if the SLACK_CHANNEL environment variable is not set
if [[ ${SLACK_CHANNEL:-"unset"} == "unset" ]]; then
  log "Environment variable SLACK_CHANNEL is not set"
  exit 1
fi

# Run driftctl scan from a hierarchical terraform state s3 repository and output it to a json file
log "Running driftctl scan"
/usr/local/bin/driftctl scan --quiet --from tfstate+s3://"${STATE_BUCKET_NAME}"/**/*.tfstate -o json://"${DRIFTCTL_JSON}"

# Publish the driftctl scan report summary to slack
log "Sending scan to slack"
/usr/local/bin/driftctl-slack 

# Copy the driftctl scan output to an s3 bucket
log "Copying scan to s3"
aws s3 --output json cp "${DRIFTCTL_JSON}" s3://"${SCAN_BUCKET_NAME}"/"${DRIFTCTL_JSON}"