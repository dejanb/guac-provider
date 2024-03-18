ARG BUILDPLATFORM="linux/arm64"
ARG BUILDERIMAGE="golang:1.21-bullseye"
ARG BASEIMAGE="gcr.io/distroless/static:nonroot"

FROM --platform=$BUILDPLATFORM $BUILDERIMAGE as builder

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT=""
ARG LDFLAGS

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    GOARM=${TARGETVARIANT}

WORKDIR /go/src/github.com/open-policy-agent/gatekeeper/test/externaldata/guac-provider

COPY . .

RUN go mod tidy

RUN go build -o provider provider.go

FROM $BASEIMAGE

WORKDIR /

COPY --from=builder /go/src/github.com/open-policy-agent/gatekeeper/test/externaldata/guac-provider/provider .

COPY --from=builder --chown=65532:65532 /go/src/github.com/open-policy-agent/gatekeeper/test/externaldata/guac-provider/certs/server.crt \
    /go/src/github.com/open-policy-agent/gatekeeper/test/externaldata/guac-provider/certs/server.key \
    /etc/ssl/certs/

USER 65532:65532

ENTRYPOINT ["/provider"]
