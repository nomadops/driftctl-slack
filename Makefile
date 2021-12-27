.PHONY: docs

docs:
	gomarkdoc --output README.md driftctl/driftctl.go s3/s3.go slack/slack.go 