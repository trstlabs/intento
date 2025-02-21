#!/usr/bin/make -f
VERSION := $(shell echo $(shell git describe --tags))
BUILDDIR ?= $(CURDIR)/build
build=i
COMMIT := $(shell git log -1 --format='%H')
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf:1.7.0
INTENTO_HOME=./
DOCKERNET_HOME=./dockernet
DOCKERNET_COMPOSE_FILE=$(DOCKERNET_HOME)/docker-compose.yml
LOCALINTENTO_HOME=./testutil/localintento
LOCALNET_COMPOSE_FILE=$(LOCALINTENTO_HOME)/localnet/docker-compose.yml
STATE_EXPORT_COMPOSE_FILE=$(LOCALINTENTO_HOME)/state-export/docker-compose.yml
LOCAL_TO_MAIN_COMPOSE_FILE=./scripts/local-to-mainnet/docker-compose.yml

# process build tags
LEDGER_ENABLED ?= true
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

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=intento \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=intentod \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" 

ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

.PHONY: build

###############################################################################
###                            Build & Clean                                ###
###############################################################################

build:
	which go
	mkdir -p $(BUILDDIR)/
	go build -mod=readonly $(BUILD_FLAGS) -trimpath -o $(BUILDDIR) ./...;

build-linux:
	GOOS=linux GOARCH=amd64 $(MAKE) build

install: go.sum
	go install $(BUILD_FLAGS) ./cmd/intentod

clean:
	rm -rf $(BUILDDIR)/*

###############################################################################
###                                CI                                       ###
###############################################################################

gosec:
	gosec -exclude-dir=deps -severity=high ./...

lint:
	golangci-lint run

###############################################################################
###                                Tests                                    ###
###############################################################################

test-unit:
	@go test -mod=readonly ./x/... ./app/...

test-unit-path:
	@go test -mod=readonly ./x/$(path)/...

test-cover:
	@go test -mod=readonly -race -coverprofile=coverage.out -covermode=atomic ./x/$(path)/...

test-integration-docker:
	bash $(DOCKERNET_HOME)/tests/run_all_tests.sh

test-integration-docker-all:
	@ALL_HOST_CHAINS=true bash $(DOCKERNET_HOME)/tests/run_all_tests.sh

###############################################################################
###                                DockerNet                                ###
###############################################################################

sync:
	@git submodule sync --recursive
	@git submodule update --init --recursive --depth 1

build-docker:
	@bash $(DOCKERNET_HOME)/build.sh -${build} ${BUILDDIR}

start-docker: stop-docker
	@bash $(DOCKERNET_HOME)/start_network.sh

start-docker-all: stop-docker build-docker
	@ALL_HOST_CHAINS=true bash $(DOCKERNET_HOME)/start_network.sh

clean-docker:
	@docker compose -f $(DOCKERNET_COMPOSE_FILE) stop
	@docker compose -f $(DOCKERNET_COMPOSE_FILE) down
	rm -rf $(DOCKERNET_HOME)/state
	docker image prune -a

stop-docker:
	@bash $(DOCKERNET_HOME)/pkill.sh
	docker compose -f $(DOCKERNET_COMPOSE_FILE) down

upgrade-build-old-binary:
	@DOCKERNET_HOME=$(DOCKERNET_HOME) BUILDDIR=$(BUILDDIR) bash $(DOCKERNET_HOME)/upgrades/build_old_binary.sh

submit-upgrade-immediately:
	UPGRADE_HEIGHT=150 bash $(DOCKERNET_HOME)/upgrades/submit_upgrade.sh

submit-upgrade-after-tests:
	UPGRADE_HEIGHT=500 bash $(DOCKERNET_HOME)/upgrades/submit_upgrade.sh

start-upgrade-integration-tests:
	PART=1 bash $(DOCKERNET_HOME)/tests/run_tests_upgrade.sh

finish-upgrade-integration-tests:
	PART=2 bash $(DOCKERNET_HOME)/tests/run_tests_upgrade.sh

upgrade-integration-tests-part-1: start-docker-all start-upgrade-integration-tests submit-upgrade-after-tests

setup-ics:
	UPGRADE_HEIGHT=150 bash $(DOCKERNET_HOME)/upgrades/setup_ics.sh

###############################################################################
###                              LocalNet                                   ###
###############################################################################
start-local-node: build
	@bash scripts/start_local_node.sh

###############################################################################
###                           Local to Mainnet                              ###
###############################################################################
start-local-to-main:
	bash scripts/local-to-mainnet/start.sh

stop-local-to-main:
	docker compose -f $(LOCAL_TO_MAIN_COMPOSE_FILE) down

###############################################################################
###                                Protobuf                                 ###
###############################################################################

containerProtoVer=0.15.1
containerProtoImage=ghcr.io/cosmos/proto-builder:$(containerProtoVer)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(DOCKER) run --user $(id -u):$(id -g) --rm -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protocgen.sh; 

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

#? proto-update-deps: Update protobuf dependencies
proto-update-deps:
	@echo "Updating Protobuf dependencies"
	$(DOCKER) run --rm -v $(CURDIR)/proto:/workspace  --user $(id -u):$(id -g) --workdir /workspace $(containerProtoImage) buf mod update 

###############################################################################
###                             Interchaintest                              ###
###############################################################################

get-heighliner:
	git clone https://github.com/strangelove-ventures/heighliner.git
	cd heighliner && go install

local-image:
ifeq (,$(shell which heighliner))
	echo 'heighliner' binary not found. Consider running `make get-heighliner`
else
	heighliner build -c intento --local --dockerfile cosmos --build-target "make install" --binaries "/go/bin/intentod"  --build-env "CGO_ENABLED=1 BUILD_TAGS=muslc"
endif


e2e-test:
	cd e2e && go test -race -v -timeout 30m
