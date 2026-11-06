ifneq (,$(wildcard ./secrets.env))
    include secrets.env
    export
endif

.PHONY: build
build:
	go build ./...

.PHONY: check
check:
ifeq ($(OS),Windows_NT)
	go test ./...
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	COVER_THRESHOLD=$(if $(and $(MOOV_PUBLIC_KEY),$(MOOV_SECRET_KEY)),30.0,0) GOTEST_PKGS=./... ./lint-project.sh
endif

.PHONY: clean
clean:
	@rm -rf ./bin/ ./tmp/ coverage.txt misspell* staticcheck lint-project.sh

.PHONY: cover-test cover-web
cover-test:
	go test -coverprofile=cover.out ./...

cover-web:
	go tool cover -html=cover.out

examples-e2ee:
	go test ./examples/e2ee/...