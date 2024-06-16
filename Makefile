
.PHONY: test_cover
test_cover:
	@go test -timeout=1s -v -race -cover -coverprofile=out.out ./...

.PHONY: test
test:test_cover
	@echo "\n"
	@go tool cover -func=out.out

.PHONY: test_html
test_html:test_cover
	@echo "\n"
	@go tool cover -html=out.out

.PHONY: test_bench
test_bench:
	@go test -bench=. -benchmem ./

.PHONY: rm_test
rm_test:
	@rm -f *.out
	@rm -f *.test

.PHONY: test_profile
test_profile:
	@go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out xerror_test.go
	@go tool pprof -http=":8081" profile.out

.PHONY: build
build:
	@go build -o main cmd/main.go
	./main xtest.go

.PHONY: protobuf
protobuf:
	protobuild vendor
	protobuild gen

vet:
	go vet ./...

refactor:
	gofumpt -l -w -extra .

lint:
	@golangci-lint run --skip-dirs-use-default --timeout 3m0s
