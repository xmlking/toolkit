ARG GO_VERSION=1
ARG UPX_VERSION=3.96

FROM --platform=$BUILDPLATFORM pratikimprowise/upx:${UPX_VERSION} AS upx
FROM golang:${GO_VERSION}-alpine AS base

COPY --from=upx /usr/local/bin/upx /bin/upx

ARG GRPC_HEALTH_PROBE_VERSION=v0.4.8
ARG GRPCURL_VERSION=1.8.6
ARG GRPCURL_SHA256=5d6768248ea75b30fba09c09ff8ba91fbc0dd1a33361b847cdaf4825b1b514a7

RUN apk --no-cache add make git gcc libtool musl-dev dumb-init \
    && rm -rf /var/cache/apk/* /tmp/* \
    && wget -q -O /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 \
    && chmod +x /bin/grpc_health_probe \
    && wget -nv  -O /tmp/grpcurl.tar.gz https://github.com/fullstorydev/grpcurl/releases/download/v${GRPCURL_VERSION}/grpcurl_${GRPCURL_VERSION}_linux_x86_64.tar.gz \
    && echo "${GRPCURL_SHA256}  /tmp/grpcurl.tar.gz" | sha256sum -c - \
    && tar -xzf /tmp/grpcurl.tar.gz -C /bin/ grpcurl \
    && chmod +x /bin/grpcurl \
    && rm /tmp/grpcurl.tar.gz

# Get dependancies - will also be cached if we won't change mod/sum
WORKDIR /
COPY ./go.mod ./go.sum ./
RUN go env -w GOPROXY="https://proxy.golang.org,direct" && go mod download && rm go.mod go.sum

# Metadata params
ARG VERSION
ARG BUILD_DATE
ARG VCS_URL=toolkit
ARG VCS_REF=1
ARG VENDOR=sumo

# Metadata
LABEL org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.name="base" \
    org.opencontainers.image.description="toolkit base docker image" \
    org.opencontainers.image.url=https://github.com/xmlking/$VCS_URL \
    org.opencontainers.image.source=https://github.com/xmlking/$VCS_URL \
    org.opencontainers.image.revision=$VCS_REF \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.authors=sumanth \
    org.opencontainers.image.vendor=$VENDOR \
    org.opencontainers.image.ref.name=$VCS_REF \
    org.opencontainers.image.licenses=MIT \
