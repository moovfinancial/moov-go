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
	COVER_THRESHOLD=35.0 GOTEST_PKGS=./... ./lint-project.sh
endif

.PHONY: setup
setup:
	docker-compose up -d --force-recreate --remove-orphans

.PHONY: teardown
teardown:
	-docker-compose down --remove-orphans
	-docker-compose rm -f -v

.PHONY: clean
clean:
	@rm -rf ./bin/ ./tmp/ coverage.txt misspell* staticcheck lint-project.sh

.PHONY: cover-test cover-web
cover-test:
	go test -coverprofile=cover.out ./...
cover-web:
	go tool cover -html=cover.out
