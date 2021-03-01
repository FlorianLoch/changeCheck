default: build

.PHONY: build test run clean docker-build docker-run dokku-deploy

change-check_bin = ./change-check
all_go_files = $(shell find . -type f -name '*.go')
all_files = $(shell find . -path $(change-check_bin) -prune -false -o -path ./.make -prune -false -o -path ./.change_check_cache -prune -false -o -type f -name '*')

git_version = $(shell git log -1 --format=%aI)
git_author_date = $(shell git describe --always)
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
	rm $(change-check_bin)

$(change-check_bin): $(all_go_files)
	go build -ldflags "-X main.gitVersion=$(git_version) -X main.gitAuthorDate=$(git_author_date) -X main.buildDate=$(build_date)" -o $(change-check_bin) .

docker-build: .make/docker-build

.make/docker-build: $(all_files)
	docker build --build-arg GIT_VERSION=$(git_version) --build-arg GIT_AUTHOR_DATE=$(git_author_date) --build-arg BUILD_DATE=$(build_date) -t fdloch/change-check .
	mkdir -p .make/ && touch .make/docker-build

docker-run: .make/docker-build
	mkdir -p .change_check_cache
	docker run --mount src="$(shell pwd)/.change_check_cache",dst=/app/.change_check_cache,type=bind --mount src="$(shell pwd)/change-check.config.yaml",dst=/app/change-check.config.yaml,type=bind --env PORT=8080 --env INTERFACE=0.0.0.0 --env APP_BASE_URL=http://localhost:8080 -p 8080:8080 fdloch/change-check

dokku-deploy: .make/dokku-deploy

.make/dokku-deploy: test .make/docker-build
	docker tag fdloch/change-check:latest dokku/change-check:latest
	docker save dokku/change-check:latest | ssh florian@vps.fdlo.ch "docker load"
	ssh -t florian@vps.fdlo.ch "sudo dokku tags:deploy change-check latest"
	touch .make/dokku-deploy