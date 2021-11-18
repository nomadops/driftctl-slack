FROM alpine:3.14

ARG OS="linux"
ARG ARCH="amd64"

WORKDIR /app
COPY driftctl-slack /usr/local/bin/driftctl-slack
RUN chmod +x /usr/local/bin/driftctl-slack
COPY driftctl-slack.sh /usr/local/bin/driftctl-slack.sh
RUN apk add curl
RUN chmod +x /usr/local/bin/driftctl-slack.sh
RUN curl -L https://github.com/cloudskiff/driftctl/releases/latest/download/driftctl_linux_amd64 -o /usr/local/bin/driftctl
RUN chmod +x /usr/local/bin/driftctl

CMD /usr/local/bin/driftctl-slack.sh