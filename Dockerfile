FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

# CGO_ENABLED=0 produces a static binary, with no libc dependency.
# -trimpath and -ldflags "-s -w" reduce the binary size.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /bin/http-server-projeto-korp .

FROM alpine:3.22

RUN addgroup -S korp && adduser -S -G korp korp
USER korp

COPY --from=build /bin/http-server-projeto-korp /usr/local/bin/http-server-projeto-korp

EXPOSE 8080

# /healthz is not instrumented, so probes don't show up in the request metrics.
# --start-interval probes every 2s during start-up, so dependent containers
# (nginx gates on service_healthy) come up seconds after the server does.
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --start-interval=2s --retries=3 \
    CMD wget -q -O /dev/null "http://127.0.0.1:${PORT:-8080}/healthz" || exit 1

ENTRYPOINT ["http-server-projeto-korp"]
