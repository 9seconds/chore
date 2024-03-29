APP_NAME := chore
GOBIN    := .bin

GOLANGCI_LINT_VERSION := v1.54.2
GOFUMPT_VERSION       := v0.5.0
GORELEASER_VERSION    := v1.20.0

STATIC_FLAGS := -buildmode=pie -modcacherw -trimpath -mod=readonly -ldflags=-linkmode=external -ldflags=-buildid='' -ldflags="-s -w"
GOTOOL       := env "GOBIN=$(abspath $(GOBIN))" "PATH=$(abspath $(GOBIN)):$(PATH)"
GO_FILES     := $(shell find . -name "*.go" -type f | grep -vE '_test\.go$$')

GPG_KEY := C9E3D1D5

# -----------------------------------------------------------------------------

.PHONY: all
all: $(APP_NAME)

$(APP_NAME): $(GO_FILES) go.sum
	@go build -tags timetzdata -o "$(APP_NAME)"

vendor: go.mod go.sum
	@$(MOD_ON) go mod vendor

$(GOBIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)/golangci-lint:
	@env GOBIN=$(abspath $(dir $@)) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

$(GOBIN)/goreleaser-$(GORELEASER_VERSION)/goreleaser:
	@env GOBIN=$(abspath $(dir $@)) go install github.com/goreleaser/goreleaser@$(GORELEASER_VERSION)

$(GOBIN)/golangci-lint: $(GOBIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)/golangci-lint
	@ln -sf $(abspath $<) $@

$(GOBIN)/gofumpt-$(GOFUMPT_VERSION)/gofumpt:
	@env GOBIN=$(abspath $(dir $@)) go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)

$(GOBIN)/gofumpt: $(GOBIN)/gofumpt-$(GOFUMPT_VERSION)/gofumpt
	@ln -sf $(abspath $<) $@

$(GOBIN)/goreleaser: $(GOBIN)/goreleaser-$(GORELEASER_VERSION)/goreleaser
	@ln -sf $(abspath $<) $@

# -----------------------------------------------------------------------------

.PHONY: static
static:
	@env go build \
		$(STATIC_FLAGS) \
		-tags netgo \
		-tags timetzdata \
		-a \
		-o "$(APP_NAME)"

.PHONY: test
test:
	@go test -parallel $(shell nproc) ./...

.PHONY: benchmark
benchmark:
	@go test -run XXXXXXX -bench=. -benchmem ./...

.PHONY: full-test
full-test:
	@go test -race -cover -coverprofile coverage.out ./...

.PHONY: lint
lint: $(GOBIN)/golangci-lint
	@$(GOTOOL) golangci-lint run ./...

.PHONY: fmt
fmt: $(GOBIN)/gofumpt
	@$(GOTOOL) gofumpt -extra -w .

.PHONY: clean
clean:
	@git clean -xfd

.PHONY: update
update:
	@go get -u && go mod tidy -go=1.21

.PHONY: clean-dist
clean-dist:
	@rm -rf dist

.PHONY: snapshot
snapshot: $(GOBIN)/goreleaser clean-dist
	@$(GOTOOL) env "GPG_TTY=$(shell tty)" "GPG_KEY=$(GPG_KEY)" goreleaser release --snapshot --clean && \
		rm dist/config.yaml && \
		find ./dist -type d | grep -vP 'dist$$' | xargs rm -r

.PHONY: release
release: $(GOBIN)/goreleaser clean-dist
	@$(GOTOOL) env "GPG_TTY=$(shell tty)" "GPG_KEY=$(GPG_KEY)" goreleaser release --skip-publish --clean && \
		rm dist/config.yaml && \
		find ./dist -type d | grep -vP 'dist$$' | xargs rm -r
