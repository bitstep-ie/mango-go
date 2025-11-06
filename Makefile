.PHONY: all
all: check build test gremlins

.PHONY: check
check: checkfmt vet

.PHONY: gremlins
gremlins: gremlins-install gremlins-run

.PHONY: gremlins-install
gremlins-install:
	go install github.com/go-gremlins/gremlins/cmd/gremlins@v0.5.0

.PHONY: gremlins-run
gremlins-run:
	gremlins unleash --config=.gremlins.yaml --output=gremlins.json

.PHONY: test
test:
	go test -v --cover ./... -coverprofile=coverage.out -covermode=count -json &> report.json

.PHONY: build
build:
	go build ./...

.PHONY: checkfmt
checkfmt:
	./scripts/gofmt.sh $(GO)

.PHONY: vet
vet:
	go vet ./...