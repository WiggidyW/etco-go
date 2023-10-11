ARG ESI_USER_AGENT='etco-go-updater default-user-agent'
ARG BUCKET_CREDS_JSON

# builder
FROM golang:1.21.3 AS builder

COPY etco-go-bucket/ /root/etco-go-bucket/
COPY etco-go-updater/ /root/etco-go-updater/
WORKDIR /root/etco-go-updater

RUN go build -o /root/out/bin .

# binary container
FROM alpine:3.18.4

RUN apk add --no-cache ca-certificates
RUN apk add --no-cache libc6-compat

COPY --from=builder /root/out/bin /root/bin
WORKDIR /root

CMD ["/root/bin"]
