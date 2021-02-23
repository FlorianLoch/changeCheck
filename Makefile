default: build

.PHONY: build test run

changeCheck_bin = ./changeCheck
all_go_files = $(shell find . -type f -name '*.go')

build: $(changeCheck_bin)

test:
ifeq (, $(shell which richgo))
	go test ./...
else
	richgo test ./...
endif

run: $(changeCheck_bin)
	$(changeCheck_bin)

$(changeCheck_bin): $(all_go_files)
	go build -o $(changeCheck_bin) .