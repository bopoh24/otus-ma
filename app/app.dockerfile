# build stage ==================================================================
FROM golang:1.21.3 as builder
WORKDIR /src
COPY .. .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./app ./cmd/app

# final stage ==================================================================
FROM istio/distroless
COPY --from=builder /src/app /app
EXPOSE 8000
ENTRYPOINT ["/app"]



