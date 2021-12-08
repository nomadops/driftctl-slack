#!/bin/bash

# PID of current shell script
PID=$$

# Name of current shell script
CMD=$(basename "${0}")

# Formats basic syslog messages as JSON 
log() {
  jq --arg message "${1}" --arg hostname "${HOSTNAME}" --arg service "${CMD}[${PID}]" -nc '{"@timestamp": now|strftime("%Y-%m-%dT%H:%M:%S%z"),"hostname": $hostname, "service": $service, "message": $message}'
}

# Check for BUCKET_NAME environment variable
if [[ ${BUCKET_NAME:-"unset"} == "unset" ]]; then
  log "BUCKET_NAME is not set"
  exit 1
fi

# Check for SCAN_BUCKET_NAME environment variable
if [[ ${SCAN_BUCKET_NAME:-"unset"} == "unset" ]]; then
  log "SCAN_BUCKET_NAME is not set"
  exit 1
fi

# Check for DRIFTCTL_JSON environment variable
if [[ ${DRIFTCTL_JSON:-"unset"} == "unset" ]]; then
  log "BUCKET_NAME is not set"
  exit 1
fi

# Run a driftctl scan
log "Running driftctl scan"
/usr/local/bin/driftctl scan --quiet --from tfstate+s3://"${BUCKET_NAME}"/**/*.tfstate -o json://"${DRIFTCTL_JSON}"

# Publish driftsctl scan summary to slack
log "Sending driftctl scan summary to slack"
/usr/local/bin/driftctl-slack 

# Copy driftctl scan summary to S3
log "Copying ${DRIFTCTL_JSON} to s3 ${SCAN_BUCKET_NAME}"
aws s3 cp "${DRIFTCTL_JSON}" s3://"${SCAN_BUCKET_NAME}"/"${DRIFTCTL_JSON}"
