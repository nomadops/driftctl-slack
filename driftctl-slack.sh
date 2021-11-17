#!/usr/bin/env bash

/usr/local/bin/driftctl scan --from tfstate+s3://"${BUCKET_NAME}"/**/*.tfstate -o console://,json://"${DRIFTCTL_JSON}" || exit 1
/usr/local/bin/driftctl-slack 