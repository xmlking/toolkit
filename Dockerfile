FROM golang:1.14-alpine

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group && \
    apk --no-cache add make git gcc libtool musl-dev ca-certificates dumb-init && \
    rm -rf /var/cache/apk/* /tmp/*  && \
    GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
    wget -q -O /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

ENV GO111MODULE=on

# Get dependancies - will also be cached if we won't change mod/sum
WORKDIR /
COPY ./go.mod ./go.sum ./
RUN go env -w GO111MODULE=on && unset GOPROXY && go env -w GOPROXY="https://proxy.golang.org,direct" && go mod download && \
    go get github.com/ahmetb/govvv && \
    go get github.com/markbates/pkger/cmd/pkger && \
    rm go.mod go.sum

# Metadata params
ARG BUILD_DATE
ARG VCS_URL=toolkit
ARG VCS_REF=1
ARG VENDOR=sumo

# Metadata
LABEL org.label-schema.build-date=$BUILD_DATE \
    org.label-schema.name="base" \
    org.label-schema.description="toolkit base docker image" \
    org.label-schema.url="https://example.com" \
    org.label-schema.vcs-url=https://github.com/xmlking/$VCS_URL \
    org.label-schema.vcs-ref=$VCS_REF \
    org.label-schema.vendor=$VENDOR \
    org.label-schema.version=$VERSION \
    org.label-schema.docker.schema-version="1.0"
