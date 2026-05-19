GO = go
GOLANGCI_LINT = $(GO) tool golangci-lint
GORELEASER = $(GO) run github.com/goreleaser/goreleaser/v2@latest

.PHONY: go-mod-tidy
go-mod-tidy:
	@echo "go mod tidy in all modules" && \
		$(GO) mod tidy -compat=1.25.0

.PHONY: lint
lint: go-mod-tidy
	@echo "Starting linting..." && \
		$(GOLANGCI_LINT) run --concurrency=4 --allow-serial-runners $(ARGS)
lint-fix: ARGS=--fix
lint-fix: lint
	@echo "✅ Lint fixing completed"

.PHONY: test
test:
	go test ./... -race
	@echo "✅ Testing completed"

.PHONY: release-check
release-check:
	@echo "Checking GoReleaser configuration..." && \
		$(GORELEASER) check
	@echo "✅ GoReleaser configuration check completed"

.PHONY: release-snapshot
release-snapshot: release-check
	@echo "Building snapshot release artifacts..." && \
		$(GORELEASER) release --snapshot --clean --skip=publish
	@echo "✅ Snapshot release build completed"

.PHONY: confirm-release
confirm-release:
	@tag=$$(git describe --tags --exact-match HEAD 2>/dev/null) || { \
	  echo; \
	  echo "Current HEAD is not tagged. Create or checkout a release tag first."; \
	  echo; \
	  exit 1; \
	}; \
	if [ "$(CONFIRM_RELEASE)" != "$$tag" ]; then \
	  echo; \
	  echo "Refusing to publish release $$tag without explicit confirmation."; \
	  echo "Run: CONFIRM_RELEASE=$$tag make release"; \
	  echo; \
	  exit 1; \
	fi; \
	echo "Confirmed release $$tag"

.PHONY: release
release: confirm-release test check-clean-work
	@echo "Building and publishing release artifacts..." && \
		$(GORELEASER) release --clean
	@echo "✅ Release completed"

.PHONY: check-clean-work
check-clean-work:
	@if ! git diff --quiet; then \
	  echo; \
	  echo 'Working tree is not clean, did you forget to run "git stash" or "git commit"?'; \
	  echo; \
	  git status; \
	  exit 1; \
	fi
