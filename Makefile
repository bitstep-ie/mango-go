# Purpose of the makefile is to replicate the CI locally

.PHONY: all
all: check build validate

# ---------------------------
# 🧹 Check Job
# - golangci-lint
# - go vet
# ---------------------------
.PHONY: check
check: lint vet

.PHONY: lint
lint:
	docker run -t --rm -v ./:/app -w /app golangci/golangci-lint:v2.6.1 golangci-lint run

.PHONY: vet
vet:
	go vet ./...

# ---------------------------
# 🏗️ Build Job
# ---------------------------
.PHONY: build
build:
	go build ./...


# ---------------------------
# 🧪 Validate Job
# ---------------------------
.PHONY: validate
validate: test gremlins

.PHONY: test
test:
	go test -v --cover ./... -coverprofile=coverage.out -covermode=count -json &> report.json

.PHONY: gremlins
gremlins: gremlins-install gremlins-run

.PHONY: gremlins-install
gremlins-install:
	go install github.com/go-gremlins/gremlins/cmd/gremlins@v0.5.0

.PHONY: gremlins-run
gremlins-run:
	gremlins unleash --config=.gremlins.yaml --output=gremlins.json


