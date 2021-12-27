<!-- Code generated by gomarkdoc. DO NOT EDIT -->

[![GitHub Actions](https://github.com/nomadops/driftctl-slack/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/nomadops/driftctl-slack/actions/workflows/ci.yml)

driftctl-slack is a golang wrapper to send driftctl scan summary to slack.

We've provided a [Terraform module](https://github.com/nomadops/terraform-aws-driftctl-slack) for deployment as a scheduleds ECS Fargate task. 
`driftctl-slack` will execute driftctl-scan, send a summary of the report to a slack channel and then copy the report to a S3 bucket.

`driftctl scan` will be executed as the following command:

```bash
$ /usr/local/bin/driftctl scan --quiet --from tfstate+s3://"${STATE_BUCKET}"/**/*.tfstate -o json://"${DRIFTCTL_JSON}"
````

Usage:
```bash
$ export SCAN_BUCKET=my-scan-bucket \
  STATE_BUCKET=my-tf-state-bucket \
  TOKEN=my-slack-token \
  CHANNEL="#gitops"
$ driftctl-slack
```

```json
{"level":"info","service":"driftctl","time":"2021-12-27T14:49:17Z"}
{"level":"info","service":"driftctl","time":"2021-12-27T14:49:17Z","message":"Driftctl scan detected drift."}
{"level":"info","service":"driftctl","time":"2021-12-27T14:49:17Z"}
{"level":"info","service":"driftctl-slack","channel":"#gitops","time":"2021-12-27T14:49:17Z","message":"Message successfully sent to slack."}
{"level":"info","service":"driftctl-slack","total_resources":428,"total_changed":0,"total_unmanaged":307,"total_missing":7,"total_managed":114,"time":"2021-12-27T14:49:17Z","message":"Driftctl scan summary"}
{"level":"info","service":"driftctl-slack","VersionId":"2AgtP5l6bRYGW30DJtT_89K_GueXeW7m","ti
```

driftctl-slack is licensed under Apache License 2.0.


# driftctl\-slack

```go
import "github.com/nomadops/driftctl-slack"
```

## Index





Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
