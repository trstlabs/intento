#!/usr/bin/make -f

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --exact-match 2>/dev/null)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
LEDGER_ENABLED ?= false
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
TM_VERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::') # grab everything after the space in "github.com/cometbft/cometbft v0.34.7"
DOCKER := $(shell which docker)

DOCKERNET_HOME=./dockernet
DOCKERNET_COMPOSE_FILE=$(DOCKERNET_HOME)/docker-compose.yml
DOCKER_IMAGE_BUILD=t
BUILDDIR ?= $(CURDIR)/build

DOCKER_TAG ?= latest

GO_SYSTEM_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1-2)
REQUIRE_GO_VERSION = 1.20

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(INTO_BUILD_OPTIONS)))
  build_tags += gcc cleveldb
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace := $(whitespace) $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=into \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=intentod \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
			-X github.com/cometbft/cometbft/version.TMCoreSemVer=$(TM_VERSION)

ifeq (cleveldb,$(findstring cleveldb,$(INTO_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq (,$(findstring nostrip,$(INTO_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(INTO_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

#$(info $$BUILD_FLAGS is [$(BUILD_FLAGS)])

# # The below include contains the tools target.
# include contrib/devtools/Makefile

###############################################################################
###                              Build                                      ###
###############################################################################

check_version:
ifneq ($(GO_SYSTEM_VERSION), $(REQUIRE_GO_VERSION))
	@echo "ERROR: Go version 1.20 is required for $(VERSION) of Intento."
	exit 1
endif

all: install lint run-tests test-e2e vulncheck

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/

$(BUILD_TARGETS): check_version go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./cmd/intentod

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

vulncheck: $(BUILDDIR)/
	GOBIN=$(BUILDDIR) go install golang.org/x/vuln/cmd/govulncheck@latest
	$(BUILDDIR)/govulncheck ./...

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go install github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd -d 2 | dot -Tpng -o dependency-graph.png

clean:
	rm -rf $(BUILDDIR)/ artifacts/

distclean: clean
	rm -rf vendor/

###############################################################################
###                                Release                                  ###
###############################################################################

# create tag and run goreleaser without publishing
create-release-dry-run:
ifneq ($(strip $(TAG)),)
	@echo "--> Dry running release for tag: $(TAG)"
	@echo "--> Create tag: $(TAG) dry run"
	git tag -s $(TAG) -m $(TAG)
	git push origin $(TAG) --dry-run
	@echo "--> Delete local tag: $(TAG)"
	@git tag -d $(TAG)
	@echo "--> Running goreleaser"
	@go install github.com/goreleaser/goreleaser@latest
	TM_VERSION=$(TM_VERSION) goreleaser release --snapshot --clean
	@rm -rf dist/
	@echo "--> Done create-release-dry-run for tag: $(TAG)"
else
	@echo "--> No tag specified, skipping tag release"
endif

# create tag and publish it
create-release:
ifneq ($(strip $(TAG)),)
	@echo "--> Running release for tag: $(TAG)"
	@echo "--> Create release tag: $(TAG)"
	git tag -s $(TAG) -m $(TAG)
	git push origin $(TAG)
	@echo "--> Done creating release tag: $(TAG)"
else
	@echo "--> No tag specified, skipping create-release"
endif


###############################################################################
###                      		     Tests		                            ###
###############################################################################

mocks: $(MOCKS_DIR)
	@go install github.com/golang/mock/mockgen@v1.6.0
	sh ./scripts/mockgen.sh
.PHONY: mocks


test-unit:
	@go test -mod=readonly ./x/alloc... ./x/mint... ./x/claim... ./x/intent... ./app/...

test-unit-path:
	@go test -mod=readonly ./x/$(path)/...

test-cover:
	@go test -mod=readonly -race -coverprofile=coverage.out -covermode=atomic ./x/$(path)/...


###############################################################################
###                                Docker                                  ###
###############################################################################

RUNNER_BASE_IMAGE_DISTROLESS := gcr.io/distroless/static-debian11
RUNNER_BASE_IMAGE_ALPINE := alpine:3.17
RUNNER_BASE_IMAGE_NONROOT := gcr.io/distroless/static-debian11:nonroot

docker-build:
	@DOCKER_BUILDKIT=1 docker build \
		-t ghcr.io/trstlabs/localintento:${DOCKER_TAG} \
		-t ghcr.io/trstlabs/localintento-distroless:${DOCKER_TAG} \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg RUNNER_IMAGE=$(RUNNER_BASE_IMAGE_DISTROLESS) \
		--build-arg GIT_VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(COMMIT) \
		-f Dockerfile .



###############################################################################
###                                Linting                                  ###
###############################################################################
golangci_lint_cmd=golangci-lint
golangci_version=v1.52.2

lint:
	@echo "--> Running linter"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run --timeout=10m

lint-fix:
	@echo "--> Running linter"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run --fix --out-format=tab --issues-exit-code=0

format:
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -path "./tests/mocks/*" -not -name "*.pb.go" -not -name "*.pb.gw.go" -not -name "*.pulsar.go" -not -path "./crypto/keys/secp256k1/*" | xargs gofumpt -w -l
	$(golangci_lint_cmd) run --fix
.PHONY: format

###############################################################################
###                                Localnet                                 ###
###############################################################################

start-localnet-ci: build
	rm -rf ~/.intentod-liveness
	./build/intentod init liveness --chain-id liveness --home ~/.intentod-liveness
	./build/intentod config chain-id liveness --home ~/.intentod-liveness
	./build/intentod config keyring-backend test --home ~/.intentod-liveness
	./build/intentod keys add val --home ~/.intentod-liveness
	./build/intentod add-genesis-account val 10000000000000000000000000stake --home ~/.intentod-liveness --keyring-backend test
	./build/intentod gentx val 1000000000stake --home ~/.intentod-liveness --chain-id liveness
	./build/intentod collect-gentxs --home ~/.intentod-liveness
	sed -i.bak'' 's/minimum-gas-prices = ""/minimum-gas-prices = "0uatom"/' ~/.intentod-liveness/config/app.toml
	./build/intentod start --home ~/.intentod-liveness --x-crisis-skip-assert-invariants

.PHONY: start-localnet-ci

###############################################################################
###                     Local Hermes+TransferPort+ICA+AutoIBCT              ###
###############################################################################


build-hermes:
	docker build -f deployment/ibc/hermes/hermes.Dockerfile -t hermes:v0.0.0 deployment/ibc/hermes

# builds a local IBC connection and channel with hermes and docker
run-localibc: build-hermes
	docker compose -f deployment/ibc/hermes/docker-compose.yml up

kill-localibc:
	docker compose -f deployment/ibc/hermes/docker-compose.yml stop 
	docker compose -f deployment/ibc/hermes/docker-compose.yml rm -f

###############################################################################
###                 		    Local Go-Relayer			                ###
###############################################################################


# runs 2 local chains
run-localchains: 
	@echo "Initializing both blockchains..."
	docker compose -f deployment/ibc/relayer/docker-compose.yml up

kill-localchains:
	@echo "Killing both blockchains..."
	docker compose -f deployment/ibc/relayer/docker-compose.yml stop 
	docker compose -f deployment/ibc/relayer/docker-compose.yml rm -f

init-golang-rly:
	@echo "Initializing relayer..."
	./deployment/ibc/relayer/interchain-acc-config/rly-init.sh

# creates a relayer for 2 local chains and adds them to rly config
create-rly: kill-dev
	@echo "Initializing relayer..."
	./deployment/ibc/relayer/init.sh

# adds juno to the rly config
create-rly-juno:
	@echo "Initializing relayer..."
	./deployment/ibc/relayer/init-juno.sh

restart-rly: @echo "Restarting relayer..."
	rly tx connection intentodev1-intentodev2 --override
	rly start intentodev1-intentodev2 -p events -b 100 --debug > rly.log

kill-dev:
	@echo "Killing intentod and removing previous data"
	-@rm -rf ./data
	-@killall intentod 2>/dev/null

# starts Into relayer given the localchains are running
run-rly: 
	@echo "Starting up local test relayer..."
	./deployment/ibc/start.sh

# starts a Into relayer, localintento 1&2 and localjuno
run-rly-juno: 
	@echo "Starting up local test relayers..."
	docker compose -f deployment/ibc/docker-compose-rly.yml up

# stops  a Into relayer, localintento 1&2 and localjuno
kill-rly-juno: 
	@echo "Stopping and removing local test relayers..."
	docker compose -f deployment/ibc/docker-compose-rly.yml stop 
	docker compose -f deployment/ibc/docker-compose-rly.yml rm -f

run-go-rly:
	./deployment/ibc/relayer/interchain-acc-config/rly-start.sh

run-localchains-juno: build-hermes
	docker compose -f deployment/ibc/relayer/docker-compose-juno.yml up

kill-localchains-juno:
	docker compose -f deployment/ibc/relayer/docker-compose-juno.yml stop 
	docker compose -f deployment/ibc/relayer/docker-compose-juno.yml rm -f

###############################################################################
###                                Swagger                                  ###
###############################################################################

# Install the runsim binary with a temporary workaround of entering an outside
# directory as the "go get" command ignores the -mod option and will polute the
# go.{mod, sum} files.
#
# ref: https://github.com/golang/go/issues/30515
statik:
	@echo "Installing statik..."
	@(cd /tmp && GO111MODULE=on go get github.com/rakyll/statik@v0.1.6)


update-swagger-docs: statik
	statik -src=client/docs/static/swagger/ -dest=client/docs -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
        exit 1;\
    else \
        echo "\033[92mSwagger docs are in sync\033[0m";\
    fi

.PHONY: update-swagger-docs statik

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.11.0
containerProtoImage=ghcr.io/cosmos/proto-builder:$(protoVer)
containerProtoGen=cosmos-sdk-proto-gen-$(protoVer)
containerProtoFmt=cosmos-sdk-proto-fmt-$(protoVer)

proto-gen:
	@echo "Generating Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGen}$$"; then docker start -a $(containerProtoGen); else docker run --name $(containerProtoGen) -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protocgen.sh; fi

proto-format:
	@echo "Formatting Protobuf files"
	@$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./proto -name "*.proto" -exec clang-format -i {} \;  

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protoc-swagger-gen.sh; 

proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

proto-check-breaking:
	@$(DOCKER_BUF) breaking --against $(HTTPS_GIT)#branch=main



###############################################################################
###                                DockerNet                                ###
###############################################################################

build-dockernet:
	@bash $(DOCKERNET_HOME)/build.sh -${DOCKER_IMAGE_BUILD} ${BUILDDIR}

start-dockernet: 
	@bash $(DOCKERNET_HOME)/start_network.sh

start-dockernet-all:
	@ALL_HOST_CHAINS=true bash $(DOCKERNET_HOME)/start_network.sh

clean-dockernet:
	@docker-compose -f $(DOCKERNET_COMPOSE_FILE) stop
	@docker-compose -f $(DOCKERNET_COMPOSE_FILE) down
	@bash $(DOCKERNET_HOME)/pkill.sh
	rm -rf $(DOCKERNET_HOME)/state
	docker image prune -a

stop-dockernet:
	@docker-compose -f $(DOCKERNET_COMPOSE_FILE) stop
	@docker-compose -f $(DOCKERNET_COMPOSE_FILE) down -v

test-integration-dockernet:
	bash $(DOCKERNET_HOME)/tests/run_all_tests.sh

upgrade-build-old-binary:
	@DOCKERNET_HOME=$(DOCKERNET_HOME) BUILDDIR=$(BUILDDIR) bash $(DOCKERNET_HOME)/upgrades/build_old_binary.sh

submit-upgrade-immediately:
	UPGRADE_HEIGHT=100 bash $(DOCKERNET_HOME)/upgrades/submit_upgrade.sh

submit-upgrade-after-tests:
	UPGRADE_HEIGHT=400 bash $(DOCKERNET_HOME)/upgrades/submit_upgrade.sh

start-upgrade-integration-tests:
	PART=1 bash $(DOCKERNET_HOME)/tests/run_tests_upgrade.sh

finish-upgrade-integration-tests:
	PART=2 bash $(DOCKERNET_HOME)/tests/run_tests_upgrade.sh

upgrade-integration-tests-part-1: start-dockernet-all start-upgrade-integration-tests submit-upgrade-after-tests