#!/bin/bash

echo "Echoing slack token"
echo "${SLCK_TOKEN}"
env
/usr/local/bin/driftctl scan --from tfstate+s3://"${BUCKET_NAME}"/**/*.tfstate -o json://"${DRIFTCTL_JSON}"
echo "Sending scan to slack"
/usr/local/bin/driftctl-slack 