FROM golang:1.22-alpine as builder
WORKDIR /service
COPY ../../go.sum go.mod ./
RUN go mod download
COPY ../.. /service
RUN go build -o app ./cmd/shortener/main.go

FROM alpine
WORKDIR /service
COPY --from=builder /service/app .
COPY migrations/* ./migrations/

EXPOSE 9090
ENTRYPOINT ["/service/app"]