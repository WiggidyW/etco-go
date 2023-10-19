## files builder
FROM golang:1.21.3 AS filebuilder

# Program Behavior Configuration
ARG PURCHASE_MAX_ACTIVE='3'
ARG MAKE_PURCHASE_COOLDOWN='10m'
ARG CANCEL_PURCHASE_COOLDOWN='10m'
ARG BOOTSTRAP_ADMIN_ID
ARG CORPORATION_ID
# ESI User Agent for Requests
ARG ESI_USER_AGENT='etco-go default-user-agent'
# ESI App Tokens
ARG STRUCTURE_INFO_WEB_REFRESH_TOKEN='BOOTSTRAP_UNSET'
ARG CORPORATION_WEB_REFRESH_TOKEN='BOOTSTRAP_UNSET'
# ESI App Configuration
ARG ESI_MARKETS_CLIENT_ID
ARG ESI_MARKETS_CLIENT_SECRET
ARG ESI_CORP_CLIENT_ID
ARG ESI_CORP_CLIENT_SECRET
ARG ESI_STRUCTURE_INFO_CLIENT_ID
ARG ESI_STRUCTURE_INFO_CLIENT_SECRET
ARG ESI_AUTH_CLIENT_ID
ARG ESI_AUTH_CLIENT_SECRET
# Bucket (GCP only for now) Configuration
ARG BUCKET_CREDS_JSON
# RemoteDB (Firestore only for now) Configuration
ARG REMOTEDB_PROJECT_ID
ARG REMOTEDB_CREDS_JSON
# Local Cache Configuration (256MB default)
ARG CCACHE_MAX_BYTES='268435456'
# Server Cache Configuration
ARG SCACHE_ADDRESS

ENV GOB_FILE_DIR='/root/out/gob'
ENV CONSTANTS_FILE_PATH='/root/out/buildconstants/buildconstants.go'

COPY etco-go-bucket/ /root/etco-go-bucket/
COPY etco-go-builder/ /root/etco-go-builder/
WORKDIR /root/etco-go-builder

RUN go get .
RUN go run .


## binary builder
FROM golang:1.21.3 AS builder

COPY etco-go-bucket/ /root/etco-go-bucket/
COPY etco-go/ /root/etco-go/
# Must come AFTER to overwrite buildconstants.go
COPY --from=filebuilder /root/out/buildconstants/buildconstants.go /root/etco-go/buildconstants/buildconstants.go
WORKDIR /root/etco-go

RUN go get .
RUN go build -o /root/out/bin .


## binary runner
FROM alpine:3.18.4

# GRPC Configuration
ARG GRPC_GO_LOG_SEVERITY_LEVEL='error'
ENV GRPC_GO_LOG_SEVERITY_LEVEL=${GRPC_GO_LOG_SEVERITY_LEVEL}
ARG GRPC_GO_LOG_VERBOSITY_LEVEL='0'
ENV GRPC_GO_LOG_VERBOSITY_LEVEL=${GRPC_GO_LOG_VERBOSITY_LEVEL}

RUN apk add --no-cache ca-certificates
RUN apk add --no-cache libc6-compat

COPY --from=builder /root/out/bin /root/bin
COPY --from=filebuilder /root/out/gob/ /root/
WORKDIR /root

CMD ["/root/bin"]
