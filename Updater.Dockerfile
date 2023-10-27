# builder
FROM golang:1.21.3 AS builder

COPY etco-go-bucket/ /root/etco-go-bucket/
COPY etco-go-updater/ /root/etco-go-updater/
WORKDIR /root/etco-go-updater

RUN go get .
RUN go build -o /root/out/bin .

# binary container
FROM alpine:3.18.4

ARG ESI_USER_AGENT='etco-go-updater default-user-agent'
ENV ESI_USER_AGENT=${ESI_USER_AGENT}
ARG BUCKET_NAMESPACE
ENV BUCKET_NAMESPACE=${BUCKET_NAMESPACE}
ARG BUCKET_CREDS_JSON
ENV BUCKET_CREDS_JSON=${BUCKET_CREDS_JSON}
ARG SKIP_SDE
ENV SKIP_SDE=${SKIP_SDE}
ARG SKIP_CORE
ENV SKIP_CORE=${SKIP_CORE}

RUN apk add --no-cache ca-certificates
RUN apk add --no-cache libc6-compat

COPY --from=builder /root/out/bin /root/bin
WORKDIR /root

CMD ["/root/bin"]
