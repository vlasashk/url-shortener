FROM golang:1.22-alpine as builder
WORKDIR /cronjob
COPY ../../go.sum go.mod ./
RUN go mod download
COPY ../.. /cronjob
RUN go build -o app ./cmd/cronjob/main.go

FROM alpine
WORKDIR /cronjob
COPY --from=builder /cronjob/app .

ENTRYPOINT ["/cronjob/app"]