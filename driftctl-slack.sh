#!/bin/sh

echo "Beginning scan"
/usr/local/bin/driftctl scan --from tfstate+s3://"${BUCKET_NAME}"/**/*.tfstate -o console://,json://"${DRIFTCTL_JSON}" || exit 1
echo "Sending scan to slack"
/usr/local/bin/driftctl-slack 