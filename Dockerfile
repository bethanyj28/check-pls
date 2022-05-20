FROM golang:1.18-alpine AS builder

WORKDIR /go/src/github.com/bethanyj28/check-pls
COPY . /go/src/github.com/bethanyj28/check-pls
RUN go mod vendor && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./...

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/github.com/bethanyj28/check-pls/app .
COPY --from=builder /go/src/github.com/bethanyj28/check-pls/.env .

ENTRYPOINT ["./app"]

EXPOSE 8080
