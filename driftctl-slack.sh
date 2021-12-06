#!/bin/bash

/usr/local/bin/driftctl scan --from tfstate+s3://"${BUCKET_NAME}"/**/*.tfstate -o json://"${DRIFTCTL_JSON}"
echo "Sending scan to slack"
/usr/local/bin/driftctl-slack 
echo "Copying scan to s3"
aws cp "${DRIFTCTL_JSON}" s3://"${SCAN_BUCKET_NAME}"/driftctl-scan.json