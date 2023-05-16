SHELL=/bin/bash
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))
name := $(shell head -1 $(current_dir)/go.mod|sed -e 's,^.*/,,g')

.DEFAULT_GOAL := run

depends_cmds := go
check:
	@for cmd in ${depends_cmds}; do command -v $$cmd >&/dev/null || (echo "No $$cmd command" && exit 1); done

clean:
	@for d in $(name); do if [[ -e $${d} ]]; then echo "==> Removing $${d}.." && rm -rf $${d}; fi done

run: check clean
	@LOG_LEVEL=debug go run . ./test/test.xlsx

sec:
	@gosec --color=false ./...
	@echo "[OK] Go security check was completed!"

help:
	@go run ./main.go -h

build: build-linux
build-linux:
	@make GOOS=linux GOARCH=amd64 _build
build-mac:
	@make GOOS=darwin GOARCH=arm64 _build
build-windows:
	@make GOOS=windows GOARCH=amd64 _build
build-android:
	@make GOOS=android GOARCH=arm64 _build
_build: check clean sec
	@env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-s -w"

deps:
	@go list -m all

update:
	@go get -u ./...

tidy:
	@go mod tidy

tidy-go:
	@ver=$(shell go version|awk '{print $$3}' |sed -e 's,go\(.*\)\..*,\1,g') && go mod tidy -go=$${ver}
