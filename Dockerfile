FROM golang:1.23-alpine AS build

WORKDIR /src

#the application uses only the standard library, so there is no go.sum
COPY go.mod ./
RUN go mod download

# CGO_ENABLED=0 produces a static binary, with no libc dependency.
# -trimpath and -ldflags "-s -w" reduce the binary size.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /bin/http-server-projeto-korp .

FROM alpine:3.20

RUN addgroup -S korp && adduser -S -G korp korp
USER korp

COPY --from=build /bin/http-server-projeto-korp /usr/local/bin/http-server-projeto-korp

EXPOSE 8080

# Checks the container's health by querying its own endpoint.
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q -O /dev/null http://127.0.0.1:8080/projeto-korp || exit 1

ENTRYPOINT ["http-server-projeto-korp"]
