.PHONY: run-all
run-all:
	docker compose up --force-recreate --build -d

.PHONY: down
down:
	docker compose down

.PHONY: lint
lint:
	golangci-lint run -c .golangci.yaml ./...

.PHONY: unit-test
unit-test:
	go test -race -coverprofile ./coverage/cover.out ./internal/... && \
    go tool cover -html=./coverage/cover.out -o ./coverage/cover.html && \
    open ./coverage/cover.html && \
    rm ./coverage/cover.out

.PHONY: integration-test
#https://github.com/golang/go/issues/65653
integration-test: export GOEXPERIMENT=nocoverageredesign
integration-test: .protoc-generate-client
	go test -race -tags=integration -coverprofile ./coverage/cover.out -coverpkg ./internal/... ./... && \
	go tool cover -html=./coverage/cover.out -o ./coverage/cover.html && \
	open ./coverage/cover.html && \
	rm ./coverage/cover.out