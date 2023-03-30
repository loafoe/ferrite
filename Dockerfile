FROM golang:alpine3.13 as builder
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build
COPY . .
RUN go build -o app

## Build final image
FROM alpine:3.17.3
LABEL maintainer="andy.lo-a-foe@philips.com"
RUN apk add --no-cache ca-certificates jq curl

WORKDIR /app
COPY --from=builder /build/app /app/app

EXPOSE 8080

CMD ["/app/app"]
