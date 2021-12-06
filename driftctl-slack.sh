#!/bin/bash

if [[ ${BUCKET_NAME:-"unset"} == "unset" ]]; then
  echo "BUCKET_NAME is not set"
  exit 1
fi

if [[ ${SCAN_BUCKET_NAME:-"unset"} == "unset" ]]; then
  echo "SCAN_BUCKET_NAME is not set"
  exit 1
fi

if [[ ${DRIFTCTL_JSON:-"unset"} == "unset" ]]; then
  echo "BUCKET_NAME is not set"
  exit 1
fi

/usr/local/bin/driftctl scan --quiet --from tfstate+s3://"${BUCKET_NAME}"/**/*.tfstate -o json://"${DRIFTCTL_JSON}"
echo "Sending scan to slack"
/usr/local/bin/driftctl-slack 
echo "Copying scan to s3"
aws s3 cp "${DRIFTCTL_JSON}" s3://"${SCAN_BUCKET_NAME}"/driftctl-scan.json