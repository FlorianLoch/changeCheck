default: build

.PHONY: build test run clean docker-build docker-run dokku-deploy coverage show-coverage

change-check_bin = ./change-check
cov_profile = ./coverage.out
all_go_files = $(shell find . -type f -name '*.go')
all_files = $(shell find . -path $(change-check_bin) -prune -false -o -path ./.make -prune -false -o -path ./.change_check_cache -prune -false -o -type f -name '*')

git_version = $(shell git describe --always)
git_author_date = $(shell git log -1 --format=%aI)
build_date = $(shell date +%Y-%m-%dT%H:%M:%S%z)

build: $(change-check_bin)

test:
ifeq (, $(shell which richgo))
	go test ./...
else
	richgo test ./...
endif

run: $(change-check_bin)
	$(change-check_bin)

clean:
	rm -rf .make
	rm $(cov_profile) || true
	rm $(change-check_bin)

$(change-check_bin): $(all_go_files)
	go build -ldflags "-X main.gitVersion=$(git_version) -X main.gitAuthorDate=$(git_author_date) -X main.buildDate=$(build_date)" -o $(change-check_bin) .

coverage: $(cov_profile)

$(cov_profile): $(all_go_files)
	# This workaround of grepping together a list of packages which do not solely contain test code seems to
	# be not necesarry with go 1.15.7 anymore...
	# https://github.com/golang/go/issues/27333
	go test ./... -coverpkg=$(shell go list ./... | grep -v test | tr "\n" ",") -coverprofile=$(cov_profile)

show-coverage: $(cov_profile)
	go tool cover -html=$(cov_profile)

docker-build: .make/docker-build

.make/docker-build: $(all_files)
	docker build --build-arg GIT_VERSION=$(git_version) --build-arg GIT_AUTHOR_DATE=$(git_author_date) --build-arg BUILD_DATE=$(build_date) -t fdloch/change-check .
	mkdir -p .make/ && touch .make/docker-build

docker-run: .make/docker-build
	mkdir -p .change_check_cache
	docker run --mount src="$(shell pwd)/.change_check_cache",dst=/app/.change_check_cache,type=bind --mount src="$(shell pwd)/change-check.config.yaml",dst=/app/change-check.config.yaml,type=bind --env PORT=8080 --env INTERFACE=0.0.0.0 --env APP_BASE_URL=http://localhost:8080 -p 8080:8080 fdloch/change-check

dokku-deploy: .make/dokku-deploy

.make/dokku-deploy: test $(all_files)
	-git remote add dokku dokku@vps.fdlo.ch:change-check
	git push dokku master
	touch .make/dokku-deploy