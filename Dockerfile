FROM golang:1.17-alpine3.14

ARG OS="linux"
ARG ARCH="amd64"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /usr/local/bin/driftctl-slack
RUN apk add curl
RUN curl -L https://github.com/cloudskiff/driftctl/releases/latest/download/driftctl_linux_amd64 -o /usr/local/bin/driftctl

# Clean-up and use a new container
FROM alpine:3.14
COPY --from=builder /usr/local/bin/driftctl-slack /usr/local/bin/driftctl-slack
COPY --from=builder /usr/local/bin/driftctl /usr/local/bin/driftctl
COPY driftctl-slack.sh /usr/local/bin/driftctl-slack.sh
RUN chmod +x /usr/local/bin/driftctl-slack.sh  /usr/local/bin/driftctl-slack /usr/local/bin/driftctl 
CMD /usr/local/bin/driftctl-slack.sh