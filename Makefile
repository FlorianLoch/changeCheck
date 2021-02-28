default: build

.PHONY: build test run clean

changeCheck_bin = ./changeCheck
all_go_files = $(shell find . -type f -name '*.go')

git_version = $(shell git log -1 --format=%aI)
git_author_date = $(shell git describe --always)
build_date = $(shell date +%Y-%m-%dT%H:%M:%S%z)

build: $(changeCheck_bin)

test:
ifeq (, $(shell which richgo))
	go test ./...
else
	richgo test ./...
endif

run: $(changeCheck_bin)
	$(changeCheck_bin)

clean:
	rm $(changeCheck_bin)

$(changeCheck_bin): $(all_go_files)
	go build -ldflags "-X main.gitVersion=$(git_version) -X main.gitAuthorDate=$(git_author_date) -X main.buildDate=$(build_date)" -o $(changeCheck_bin) .