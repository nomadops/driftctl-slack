// Usage:
//
// 		$ export SCAN_BUCKET=my-scan-bucket STATE_BUCKET=my-tf-state-bucket TOKEN=my-slack-token CHANNEL="#gitops"
//		$ driftctl-slack
//		{"level":"info","service":"driftctl","time":"2021-12-27T14:49:17Z"}
//		{"level":"info","service":"driftctl","time":"2021-12-27T14:49:17Z","message":"Driftctl scan detected drift."}
//		{"level":"info","service":"driftctl","time":"2021-12-27T14:49:17Z"}
//		{"level":"info","service":"driftctl-slack","channel":"#gitops","time":"2021-12-27T14:49:17Z","message":"Message successfully sent to slack."}
//		{"level":"info","service":"driftctl-slack","total_resources":428,"total_changed":0,"total_unmanaged":307,"total_missing":7,"total_managed":114,"time":"2021-12-27T14:49:17Z","message":"Driftctl scan summary"}
//		{"level":"info","service":"driftctl-slack","VersionId":"2AgtP5l6bRYGW30DJtT_89K_GueXeW7m","ti
//
package main
