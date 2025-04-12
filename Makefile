GO := go
GOFLAGS ?= -mod=readonly -ldflags "-s -w -X 'main.version=$(VERSION)-dev' -extldflags '-static'"

VERSION ?= $(shell git describe --tags --first-parent --abbrev=0 | cut -c 2-)

cover:
	$(GO) test -race -coverprofile=coverage.out ./...

cover-func: cover
	$(GO) tool cover -func=coverage.out

cover-html: cover
	$(GO) tool cover -html=coverage.out -o coverage.html

build:
	$(GO) build -a -trimpath -tags netgo $(GOFLAGS) -o bin/ ./cmd/...

fast-build:
	$(GO) build -o bin/ ./cmd/...

testclean:
	@rm -f coverage.out
	@rm -fr output

distclean:
	@rm -fr bin
	@rm -fr build
	@rm -fr dist

clean:
	@$(GO) clean

gofmt:
	@mkdir -p output
	@rm -f output/lint.log

	gofmt -d -s . 2>&1 | tee output/lint.log

	@[ ! -s output/lint.log ]

	@rm -fr output

tidy:
	@$(GO) mod tidy -v

.PHONY: build
