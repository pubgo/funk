
GO_TEST_TIMEOUT ?= 15m

.PHONY: test_cover test test_html test_bench rm_test test_profile vet refactor lint protobuf protolint changelog changelog-latest release-dry install-tools

test_cover:
	@go test -timeout=$(GO_TEST_TIMEOUT) -v -race -cover -coverprofile=out.out ./...

.PHONY: test
test: test_cover
	@echo "\n"
	@go tool cover -func=out.out

.PHONY: test_html
test_html: test_cover
	@echo "\n"
	@go tool cover -html=out.out

.PHONY: test_bench
test_bench:
	@go test -bench=. -benchmem ./...

.PHONY: rm_test
rm_test:
	@rm -f *.out
	@rm -f *.test

.PHONY: test_profile
test_profile:
	@go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out xerror_test.go
	@go tool pprof -http=":8081" profile.out

protobuf:
	protobuild vendor
	protobuild gen

protolint:
	protobuild lint

vet:
	go vet ./...

refactor:
	gofumpt -l -w -extra .

lint:
	golangci-lint run --timeout=10m --verbose

changelog:
	@command -v git-cliff >/dev/null 2>&1 || { echo "install: https://git-cliff.org/docs/installation/ (e.g. cargo install git-cliff --locked --version 2.10.1)"; exit 1; }
	git cliff --unreleased --strip header

changelog-latest:
	@command -v git-cliff >/dev/null 2>&1 || { echo "install: https://git-cliff.org/docs/installation/ (e.g. cargo install git-cliff --locked --version 2.10.1)"; exit 1; }
	git cliff --latest --strip header

release-dry:
	@command -v goreleaser >/dev/null 2>&1 || { echo "install: https://goreleaser.com/install/"; exit 1; }
	goreleaser release --snapshot --clean

install-tools:
	go install -v ./cmds/protoc-gen-go-errors2
	go install -v ./cmds/protoc-gen-go-enum2
	go install -v ./cmds/protoc-gen-go-sql2
	go install -v ./cmds/protoc-gen-go-cloudevent2
