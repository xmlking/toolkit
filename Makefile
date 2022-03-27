# Usage:
# make        	# compile all binary
# make clean  	# remove ALL binaries and objects
# make release  # add git TAG and push
GITHUB_REPO_OWNER 				:= xmlking
GITHUB_REPO_NAME 				:= toolkit
GITHUB_RELEASES_UI_URL 			:= https://github.com/$(GITHUB_REPO_OWNER)/$(GITHUB_REPO_NAME)/releases
GITHUB_RELEASES_API_URL 		:= https://api.github.com/repos/$(GITHUB_REPO_OWNER)/$(GITHUB_REPO_NAME)/releases
GITHUB_RELEASE_ASSET_URL		:= https://uploads.github.com/repos/$(GITHUB_REPO_OWNER)/$(GITHUB_REPO_NAME)/releases
GITHUB_DEPLOY_API_URL			:= https://api.github.com/repos/$(GITHUB_REPO_OWNER)/$(GITHUB_REPO_NAME)/deployments
DOCKER_REGISTRY 				:= ghcr.io
# DOCKER_REGISTRY 				:= us.gcr.io
DOCKER_CONTEXT_PATH 			:= $(GITHUB_REPO_OWNER)/$(GITHUB_REPO_NAME)
# DOCKER_REGISTRY 				:= docker.io
# DOCKER_CONTEXT_PATH 			:= xmlking
BASE_VERSION					:= latest

VERSION					:= $(shell git describe --tags || echo "HEAD")
GOPATH					:= $(shell go env GOPATH)
CODECOV_FILE 			:= build/coverage.txt
TIMEOUT  				:= 60s
# don't override
GIT_TAG					:= $(shell git describe --tags --abbrev=0 --always --match "v*")
GIT_DIRTY 				:= $(shell git status --porcelain 2> /dev/null)
GIT_BRANCH  			:= $(shell git rev-parse --abbrev-ref HEAD)
HAS_GOVVV				:= $(shell command -v govvv 2> /dev/null)
HAS_KO					:= $(shell command -v ko 2> /dev/null)
HTTPS_GIT 				:= https://github.com/$(GITHUB_REPO_OWNER)/$(GITHUB_REPO_NAME).git

# Type of service e.g api, service, web, cmd (default: "service")
TYPE = $(or $(word 2,$(subst -, ,$*)), service)
override TYPES:= service
# Target for running the action
TARGET = $(word 1,$(subst -, ,$*))

override VERSION_PACKAGE = $(shell go list ./internal/config)

# $(warning TYPES = $(TYPE), TARGET = $(TARGET))
# $(warning VERSION = $(VERSION), HAS_GOVVV = $(HAS_GOVVV), HAS_KO = $(HAS_KO))

.PHONY: all tools check_dirty sync
.PHONY: lint lint-% outdated
.PHONY: format format-%
.PHONY: release/draft release/publish

all: build

################################################################################
# Target: tools
################################################################################

tools:
	@echo "==> Installing dev tools"
	# go install github.com/ahmetb/govvv@latest
	# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	# go install github.com/bufbuild/buf/cmd/buf@latest
	# go install github.com/rvflash/goup@latest

check_dirty:
ifdef GIT_DIRTY
	$(error "Won't run on a dirty working copy. Commit or stash and try again.")
endif

################################################################################
# Target: go-mod                                                               #
################################################################################

sync:
	@echo Removing all go.sum files and download sync deps....
	@for d in `find * -name 'go.mod'`; do \
		pushd `dirname $$d` >/dev/null; \
	    echo delete and sync $$d; \
	    rm -f go.sum; \
	    go mod download; \
	    go mod tidy -compat=1.18; \
		popd >/dev/null; \
	done

verify:
	@echo Go mod verify and tidy....
	@for d in `find * -name 'go.mod'`; do \
		pushd `dirname $$d` >/dev/null; \
		echo verifying $$d; \
		go mod verify; \
		go mod tidy -compat=1.18; \
		popd >/dev/null; \
	done

outdated:
	@goup -v -m ./...

################################################################################
# Target: lints                                                                #
################################################################################

lint lint-%:
	@if [ -z $(TARGET) ]; then \
		echo "Linting all go"; \
		${GOPATH}/bin/golangci-lint run ./... --deadline=5m --config=.github/linters/.golangci.yml; \
		echo "✓ Go: Linted"; \
	else \
		echo "Linting go in ${TARGET}-${TYPE}..."; \
		${GOPATH}/bin/golangci-lint run ./${TYPE}/${TARGET}/... --config=.github/linters/.golangci.yml; \
		echo "✓ Go: Linted in ${TYPE}/${TARGET}..."; \
	fi

# @clang-format -i $(shell find . -type f -name '*.proto')

format format-%:
	@if [ -z $(TARGET) ]; then \
		echo "Formatting all go"; \
		gofmt -l -w . ; \
		echo "✓ Go: Formatted"; \
	else \
		echo "Formatting go in ${TYPE}/${TARGET}..."; \
		gofmt -l -w ./${TYPE}/${TARGET}/ ; \
		echo "✓ Go: Formatted in ${TYPE}/${TARGET}..."; \
	fi

################################################################################
# Target: tests                                                                #
################################################################################

TEST_TARGETS := test-default test-bench test-unit test-inte test-e2e test-race test-cover
.PHONY: $(TEST_TARGETS) check test tests
test-bench:   	ARGS=-run=__absolutelynothing__ -bench=. ## Run benchmarks
test-unit:   	ARGS=-short        					## Run only unit tests
test-inte:   	ARGS=-run Integration       		## Run only integration tests
test-e2e:   	ARGS=-run E2E       				## Run only E2E tests
test-race:    	ARGS=-race         					## Run tests with race detector
test-cover:   	ARGS=-cover -short -coverprofile=${CODECOV_FILE} -covermode=atomic ## Run tests in verbose mode with coverage reporting
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
check test tests:
	@if [ -z $(TARGET) ]; then \
		echo "Running $(NAME:%=% )tests for all"; \
		go test -timeout $(TIMEOUT) $(ARGS) ./... ; \
	else \
		echo "Running $(NAME:%=% )tests for ${TARGET}-${TYPE}"; \
		go test -timeout $(TIMEOUT) -v $(ARGS) ./${TYPE}/${TARGET}/... ; \
	fi

################################################################################
# Target: release                                                              #
################################################################################

release: verify
	@if [ -z $(TAG) ]; then \
		echo "no  TAG. Usage: make release TAG=v0.1.1"; \
	else \
		for m in `find * -name 'go.mod' -mindepth 1 -exec dirname {} \;`; do \
			hub release create -m "$$m/${TAG} release" $$m/${TAG}; \
		done \
	fi

release/draft: check_dirty
	@echo Publishing Draft: $(VERSION)
	@git tag -a $(VERSION) -m "[skip ci] Release: $(VERSION)" || true
	@git push origin $(VERSION)
	@echo "\n\nPlease inspect the release and run `make release/publish` if it looks good"
	@open "$(GITHUB_RELEASES_UI_URL)/$(VERSION)"

release/publish:
	@echo Publishing Release: $(VERSION)
