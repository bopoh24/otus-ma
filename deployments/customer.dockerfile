# build stage ==================================================================
FROM golang:latest as builder
WORKDIR /src
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./customer ./customer/cmd

# final stage ==================================================================
FROM istio/distroless
COPY --from=builder /src/customer /customer
EXPOSE 80
ENTRYPOINT ["/customer"]



