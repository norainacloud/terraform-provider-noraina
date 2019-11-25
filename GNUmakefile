TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -w -s $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...

tools:
	@echo "==> Installing required tooling..."
	@GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint
	@GO111MODULE=on go install github.com/bflad/tfproviderlint/cmd/tfproviderlint

.PHONY: build test fmt fmtcheck lint tools
