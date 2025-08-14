FROM alpine:3.22 as base
RUN apk add --no-cache ca-certificates
RUN adduser -D hetzner-nuke

FROM golang:1.25 AS build
COPY / /src
WORKDIR /src
ENV CGO_ENABLED=0
RUN \
  --mount=type=cache,target=/go/pkg \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -ldflags '-s -w -extldflags="-static"' -o bin/hetzner-nuke main.go

FROM base AS goreleaser
ENTRYPOINT ["/usr/local/bin/hetzner-nuke"]
COPY hetzner-nuke /usr/local/bin/hetzner-nuke
USER hetzner-nuke

FROM base
ENTRYPOINT ["/usr/local/bin/hetzner-nuke"]
COPY --from=build --chmod=755 /src/bin/hetzner-nuke /usr/local/bin/hetzner-nuke
RUN chmod +x /usr/local/bin/hetzner-nuke
USER hetzner-nuke