MODULE  := github.com/DavDaz/llm-wiki-generator
BINARY  := llm-wiki
VERSION ?= dev
LDFLAGS := -ldflags "-X '$(MODULE)/internal/version.Version=$(VERSION)'"

DEFAULT_RELEASE_BRANCH ?= main

.PHONY: build test lint vet tidy clean preflight validate-release-version validate-release-branch release

build:
	go build $(LDFLAGS) -o bin/$(BINARY) ./cmd/$(BINARY)

test:
	go test -race ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/


preflight:
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Error: working tree must be clean before release"; \
		git status --short; \
		exit 1; \
	fi
	@echo "Preflight: running go test ./..."
	go test ./...
	@echo "Preflight passed."

validate-release-version:
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "dev" ]; then \
		echo "Usage: make release VERSION=v0.3.0"; \
		exit 1; \
	fi
	@if ! echo "$(VERSION)" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$$'; then \
		echo "Error: VERSION must be semver (e.g. v0.3.0)"; \
		exit 1; \
	fi

validate-release-branch:
	@current_branch="$$(git rev-parse --abbrev-ref HEAD)"; \
		default_branch="$$(git symbolic-ref --quiet --short refs/remotes/origin/HEAD 2>/dev/null)"; \
		default_branch="$${default_branch#origin/}"; \
		if [ -z "$$default_branch" ]; then \
			default_branch="$(DEFAULT_RELEASE_BRANCH)"; \
		fi; \
		if [ "$$current_branch" != "$$default_branch" ]; then \
			echo "Error: releases must be tagged from $$default_branch (current: $$current_branch)"; \
			exit 1; \
		fi; \
		git fetch origin "$$default_branch"; \
		local_head="$$(git rev-parse HEAD)"; \
		remote_head="$$(git rev-parse "origin/$$default_branch")"; \
		if [ "$$local_head" != "$$remote_head" ]; then \
			echo "Error: local $$default_branch must match origin/$$default_branch before release"; \
			exit 1; \
		fi

release: validate-release-version validate-release-branch preflight
	@echo "Releasing $(VERSION)..."
	git tag $(VERSION)
	git push origin $(VERSION)
	@echo "Done — GitHub Actions will build and publish the release."
