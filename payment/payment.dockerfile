# build stage ==================================================================
FROM golang:latest as builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./payment/app ./payment/cmd

# final stage ==================================================================
FROM gcr.io/distroless/base-debian10
COPY --from=builder /src/payment/app /payment
EXPOSE 80
ENTRYPOINT ["/payment"]



