FROM golang:1.17-alpine3.14 AS builder

ARG OS="linux"
ARG ARCH="amd64"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build
RUN apk add curl
RUN curl -s -L https://github.com/cloudskiff/driftctl/releases/latest/download/driftctl_linux_amd64 -o driftctl.app
RUN ls -l /usr/local/bin

# Clean-up and use a new container
FROM alpine:3.15
# RUN apk add bash jq curl
COPY --from=builder /app/driftctl-slack /usr/local/bin/driftctl-slack
COPY --from=builder /app/driftctl.app /usr/local/bin/driftctl
RUN chmod +x   /usr/local/bin/driftctl-slack /usr/local/bin/driftctl 

CMD /usr/local/bin/driftctl-slack