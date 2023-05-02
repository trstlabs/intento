PACKAGES=$(shell go list ./... | grep -v '/simulation')
VERSION ?= $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
CURRENT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf

LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
BUILD_PROFILE ?= release
DEB_BIN_DIR ?= /usr/local/bin
DEB_LIB_DIR ?= /usr/lib


SGX_MODE ?= SW
BRANCH ?= develop

DOCKER_TAG ?= latest

ifeq ($(SGX_MODE), HW)
	ext := hw
else ifeq ($(SGX_MODE), SW)
	ext := sw
else
$(error SGX_MODE must be either HW or SW)
endif

CUR_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error "gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false")
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning "OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988)")
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error "gcc not installed for ledger support, please install or set LEDGER_ENABLED=false")
      else
        build_tags += ledger
      endif
    endif
  endif
endif

IAS_BUILD = sw

ifeq ($(SGX_MODE), HW)
  ifneq (,$(findstring production,$(FEATURES)))
    IAS_BUILD = production
  else
    IAS_BUILD = develop
  endif

  build_tags += hw
endif

build_tags += $(IAS_BUILD)

ifeq ($(DB_BACKEND),rocksdb)
  build_tags += gcc
endif
ifeq ($(DB_BACKEND),cleveldb)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=trst \
-X github.com/cosmos/cosmos-sdk/version.AppName=trstd \
	-X github.com/trstlabs/trst/cmd/trstd/version.ClientName=trstd \
-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags)"

ifeq ($(DB_BACKEND),cleveldb)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq ($(DB_BACKEND),rocksdb)
  CGO_ENABLED=1
  build_tags += rocksdb
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb
  ldflags += -extldflags "-lrocksdb -llz4"
endif



ldflags += -s -w
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

GO_TAGS := $(build_tags)
# -ldflags
LD_FLAGS := $(ldflags)

all: build_all

vendor:
	cargo vendor third_party/vendor --manifest-path third_party/build/Cargo.toml

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

build_cli:
	go build -o trstcli -mod=readonly -tags "$(GO_TAGS) trstcli" -ldflags '$(LD_FLAGS)' ./cmd/trstd

xgo_build_trstcli: go.sum
	@echo "--> WARNING! This builds from origin/$(CURRENT_BRANCH)!"
	xgo --targets $(XGO_TARGET) -tags="$(GO_TAGS) trstcli" -ldflags '$(LD_FLAGS)' --branch "$(CURRENT_BRANCH)" github.com/trstlabs/trst/cmd/trstd

build_local_no_rust: bin-data-$(IAS_BUILD)
	cp go-cosmwasm/target/$(BUILD_PROFILE)/libgo_cosmwasm.so go-cosmwasm/api
	go build -mod=readonly -tags "$(GO_TAGS)" -ldflags '$(LD_FLAGS)' ./cmd/trstd

build-linux: _build-linux build_local_no_rust build_cli
_build-linux: vendor
	BUILD_PROFILE=$(BUILD_PROFILE) FEATURES=$(FEATURES) FEATURES_U=$(FEATURES_U) $(MAKE) -C go-cosmwasm build-rust

build-linux-with-query: _build-linux-with-query build_local_no_rust build_cli
_build-linux-with-query: vendor
	BUILD_PROFILE=$(BUILD_PROFILE) FEATURES=$(FEATURES) FEATURES_U=query-node,$(FEATURES_U) $(MAKE) -C go-cosmwasm build-rust

build_windows_cli:
	$(MAKE) xgo_build_trstcli XGO_TARGET=windows/amd64
	mv trstd-windows-* trstcli-windows-amd64.exe

build_macos_cli:
	$(MAKE) xgo_build_trstcli XGO_TARGET=darwin/amd64
	mv trstd-darwin-* trstcli-macos-amd64

build_linux_cli:
	$(MAKE) xgo_build_trstcli XGO_TARGET=linux/amd64
	mv trstd-linux-amd64 trstcli-linux-amd64

build_linux_arm64_cli:
	$(MAKE) xgo_build_trstcli XGO_TARGET=linux/arm64
	mv trstd-linux-arm64 trstcli-linux-arm64

build_all: build-linux build_windows_cli build_macos_cli build_linux_arm64_cli

deb: build-linux deb-no-compile

#redundant
deb-query: build-linux-with-query deb-no-compile

deb-no-compile:
    ifneq ($(UNAME_S),Linux)
		exit 1
    endif
	rm -rf /tmp/trst

	mkdir -p /tmp/trst/deb/$(DEB_BIN_DIR)
	cp -f ./trstcli /tmp/trst/deb/$(DEB_BIN_DIR)/trstcli
	cp -f ./trstd /tmp/trst/deb/$(DEB_BIN_DIR)/trstd
	chmod +x /tmp/trst/deb/$(DEB_BIN_DIR)/trstd /tmp/trst/deb/$(DEB_BIN_DIR)/trstcli

	mkdir -p /tmp/trst/deb/$(DEB_LIB_DIR)
	cp -f ./go-cosmwasm/api/libgo_cosmwasm.so ./go-cosmwasm/librust_cosmwasm_enclave.signed.so /tmp/trst/deb/$(DEB_LIB_DIR)/
	chmod +x /tmp/trst/deb/$(DEB_LIB_DIR)/lib*.so

	mkdir -p /tmp/trst/deb/DEBIAN
	cp ./deployment/deb/control /tmp/trst/deb/DEBIAN/control
	printf "Version: " >> /tmp/trst/deb/DEBIAN/control
	printf "$(VERSION)" >> /tmp/trst/deb/DEBIAN/control
	echo "" >> /tmp/trst/deb/DEBIAN/control
	cp ./deployment/deb/postinst /tmp/trst/deb/DEBIAN/postinst
	chmod 755 /tmp/trst/deb/DEBIAN/postinst
	cp ./deployment/deb/postrm /tmp/trst/deb/DEBIAN/postrm
	chmod 755 /tmp/trst/deb/DEBIAN/postrm
	cp ./deployment/deb/triggers /tmp/trst/deb/DEBIAN/triggers
	chmod 755 /tmp/trst/deb/DEBIAN/triggers
	dpkg-deb --build /tmp/trst/deb/ .
	-rm -rf /tmp/trst

rename_for_release:
	-rename "s/windows-4.0-amd64/v${VERSION}-win64/" *.exe
	-rename "s/darwin-10.6-amd64/v${VERSION}-osx64/" *darwin*

sign_for_release: rename_for_release
	sha256sum trst*.deb > SHA256SUMS
	-sha256sum trstd-* trstcli-* >> SHA256SUMS
	gpg -u 91831DE812C6415123AFAA7B420BF1CB005FBCE6 --digest-algo sha256 --clearsign --yes SHA256SUMS
	rm -f SHA256SUMS

release: sign_for_release
	rm -rf ./release/
	mkdir -p ./release/
	cp trst_*.deb ./release/
	cp trstcli-* ./release/
	cp trstd-* ./release/
	cp SHA256SUMS.asc ./release/

clean:
	-rm -rf /tmp/trst
	-rm -f ./trstcli*
	-rm -f ./trstd*
	-find -name librust_cosmwasm_enclave.signed.so -delete
	-find -name libgo_cosmwasm.so -delete
	-find -name '*.so' -delete
	-find -name 'target' -type d -exec rm -rf \;
	-rm -f ./trst*.deb
	-rm -f ./SHA256SUMS*
	-rm -rf ./third_party/vendor/
	-rm -rf ./trustlesshub/.sgx_secrets/*
	-rm -rf ./x/compute/internal/keeper/trustlesshub/.sgx_secrets/*
	-rm -rf ./*.der
	-rm -rf ./x/compute/internal/keeper/*.der
	-rm -rf ./cmd/trstd/ias_bin*
	$(MAKE) -C go-cosmwasm clean-all
	$(MAKE) -C cosmwasm/enclaves/test clean

build-rocksdb-image:
	docker build --build-arg BUILD_VERSION=v6.24.3 -f deployment/dockerfiles/db-compile.Dockerfile -t trstlabs/rocksdb:v6.24.3 .

build-localtrst:_localtrst-compile
	DOCKER_BUILDKIT=1 docker build --build-arg SGX_MODE=SW --build-arg TRST_NODE_TYPE=BOOTSTRAP --build-arg CHAIN_ID=trst_chain_1 -f deployment/dockerfiles/release.Dockerfile -t build-release .
	DOCKER_BUILDKIT=1 docker build --build-arg SGX_MODE=SW --build-arg TRST_NODE_TYPE=BOOTSTRAP --build-arg CHAIN_ID=trst_chain_1 -f deployment/dockerfiles/dev-image.Dockerfile -t ghcr.io/trstlabs/localtrst:${DOCKER_TAG} .

_localtrst-compile:
	DOCKER_BUILDKIT=1 docker build \
				--build-arg BUILD_VERSION=${VERSION} \
				--build-arg FEATURES="${FEATURES},debug-print" \
				--build-arg FEATURES_U=${FEATURES_U} \
				--secret id=API_KEY,src=.env.local \
				--secret id=SPID,src=.env.local \
				--build-arg SGX_MODE=SW \
				-f deployment/dockerfiles/base.Dockerfile \
				-t rust-go-base-image \
				.


build-dev-image:_dev-trst-compile
	 docker build --build-arg SGX_MODE=SW --build-arg TRST_NODE_TYPE=BOOTSTRAP --build-arg CHAIN_ID=trst_chain_1 -f deployment/dockerfiles/release.Dockerfile -t build-release .
	 docker build --build-arg SGX_MODE=SW --build-arg TRST_NODE_TYPE=BOOTSTRAP --build-arg CHAIN_ID=trst_chain_1 -f deployment/dockerfiles/dev-image.Dockerfile -t trstlabs/trst-sw-dev:${DOCKER_TAG} .

_dev-trst-compile:

	 docker build \
				--build-arg BUILD_VERSION=${VERSION} \
				--build-arg FEATURES="${FEATURES},debug-print" \
				--build-arg FEATURES_U=${FEATURES_U} \
				--secret id=API_KEY,src=.env.dev \
				--secret id=SPID,src=.env.dev \
				--build-arg SGX_MODE=SW \
				-f deployment/dockerfiles/base.Dockerfile \
				-t rust-go-base-image \
				.

build-custom-dev-image:
    # .dockerignore excludes .so files so we rename these so that the dockerfile can find them
	cd go-cosmwasm/api && cp libgo_cosmwasm.so libgo_cosmwasm.so.x
	cd cosmwasm/enclaves/execute && cp librust_cosmwasm_enclave.signed.so librust_cosmwasm_enclave.signed.so.x
	docker build --build-arg SGX_MODE=SW --build-arg TRST_NODE_TYPE=BOOTSTRAP -f deployment/dockerfiles/custom-node.Dockerfile -t trstlabs/trst-sw-dev-custom-bootstrap:${DOCKER_TAG} .
	docker build --build-arg SGX_MODE=SW --build-arg TRST_NODE_TYPE=NODE -f deployment/dockerfiles/custom-node.Dockerfile -t trstlabs/trst-sw-dev-custom-node:${DOCKER_TAG} .
    # delete the copies created above
	rm go-cosmwasm/api/libgo_cosmwasm.so.x cosmwasm/enclaves/execute/librust_cosmwasm_enclave.signed.so.x

build-testnet: _docker_base
	@mkdir build 2>&3 || true
	docker build --build-arg BUILD_VERSION=${VERSION} --build-arg SGX_MODE=HW --build-arg TRST_NODE_TYPE=BOOTSTRAP -f deployment/dockerfiles/release.Dockerfile -t trstlabs/trst-bootstrap:v$(VERSION)-testnet .
	docker build --build-arg BUILD_VERSION=${VERSION} --build-arg SGX_MODE=HW --build-arg TRST_NODE_TYPE=NODE -f deployment/dockerfiles/release.Dockerfile -t trstlabs/trst-node:v$(VERSION)-testnet .
	docker build --build-arg SGX_MODE=HW -f deployment/dockerfiles/build-deb.Dockerfile -t deb_build .
	docker run -e VERSION=${VERSION} -v $(CUR_DIR)/build:/build deb_build

build-mainnet: _docker_base
	@mkdir build 2>&3 || true
	docker build --build-arg SGX_MODE=HW --build-arg TRST_NODE_TYPE=BOOTSTRAP -f deployment/dockerfiles/release.Dockerfile -t trstlabs/trst-bootstrap:v$(VERSION)-mainnet .
	docker build --build-arg SGX_MODE=HW --build-arg TRST_NODE_TYPE=NODE -f deployment/dockerfiles/release.Dockerfile -t trstlabs/trst-node:v$(VERSION)-mainnet .
	docker build --build-arg BUILD_VERSION=${VERSION} --build-arg SGX_MODE=HW -f deployment/dockerfiles/build-deb.Dockerfile -t deb_build .
	docker run -e VERSION=${VERSION} -v $(CUR_DIR)/build:/build deb_build


 _docker_base:
	docker build \
				--build-arg BUILD_VERSION=${VERSION} \
				--build-arg FEATURES=${FEATURES} \
				--build-arg FEATURES_U=${FEATURES_U} \
				--build-arg SGX_MODE=${SGX_MODE} \
				--secret id=API_KEY,src=api_key.txt \
				--secret id=SPID,src=spid.txt \
				-f deployment/dockerfiles/base.Dockerfile \
				-t rust-go-base-image \
				.


build-ibc-hermes:
	docker build -f deployment/dockerfiles/ibc/hermes.Dockerfile -t hermes:v0.0.0 deployment/dockerfiles/ibc

docker_bootstrap: _docker_base
	docker build --build-arg SGX_MODE=${SGX_MODE} --build-arg TRST_NODE_TYPE=BOOTSTRAP -f deployment/dockerfiles/local-node.Dockerfile -t trstlabs/trst-bootstrap-${ext}:${DOCKER_TAG} .

docker_node: _docker_base
	docker build --build-arg SGX_MODE=${SGX_MODE} --build-arg TRST_NODE_TYPE=NODE -f deployment/dockerfiles/local-node.Dockerfile -t trstlabs/trst-node-${ext}:${DOCKER_TAG} .

clean-files:
	-rm -rf /trst

	-rm -f ./trstd*
#   -find -name librust_cosmwasm_enclave.signed.so -delete
#   -find -name libgo_cosmwasm.so -delete
#   -find -name '*.so' -delete
#   -find -name 'target' -type d -exec rm -rf \;
	-rm -f ./SHA256SUMS*
	-rm -rf ./trustlesshub/.sgx_secrets/*
	-rm -rf ./x/compute/internal/keeper/trustlesshub/.sgx_secrets/*
	-rm -rf ./*.der
	-rm -rf ./x/compute/internal/keeper/*.der
	-rm -rf ./cmd/trstd/ias_bin*

docker_local_azure_hw: _docker_base
	docker build --build-arg SGX_MODE=HW --build-arg TRST_NODE_TYPE=NODE -f deployment/dockerfiles/local-node.Dockerfile -t ci-trst-sgx-node .
	docker build --build-arg SGX_MODE=HW --build-arg TRST_NODE_TYPE=BOOTSTRAP -f deployment/dockerfiles/local-node.Dockerfile -t ci-trst-sgx-bootstrap .

docker_enclave_test:
	docker build --build-arg FEATURES="test ${FEATURES}" --build-arg SGX_MODE=${SGX_MODE} -f deployment/dockerfiles/enclave-test.Dockerfile -t rust-enclave-test .

# while developing:
build-enclave: vendor
	$(MAKE) -C cosmwasm/enclaves/execute enclave

# while developing:
check-enclave:
	$(MAKE) -C cosmwasm/enclaves/execute check

# while developing:
clippy-enclave:
	$(MAKE) -C cosmwasm/enclaves/execute clippy

# while developing:
clean-enclave:
	$(MAKE) -C cosmwasm/enclaves/execute clean

sanity-test:
	SGX_MODE=SW $(MAKE) build-linux
	cp ./cosmwasm/enclaves/execute/librust_cosmwasm_enclave.signed.so .
	SGX_MODE=SW ./cosmwasm/testing/sanity-test.sh

sanity-test-hw:
	$(MAKE) build-linux
	cp ./cosmwasm/enclaves/execute/librust_cosmwasm_enclave.signed.so .
	./cosmwasm/testing/sanity-test.sh

callback-sanity-test:
	SGX_MODE=SW $(MAKE) build-linux
	cp ./cosmwasm/enclaves/execute/librust_cosmwasm_enclave.signed.so .
	SGX_MODE=SW ./cosmwasm/testing/callback-test.sh

build-test-contracts:
	# echo "" | sudo add-apt-repository ppa:hnakamur/binaryen
	# sudo apt update
	# sudo apt install -y binaryen
	$(MAKE) -C ./x/compute/internal/keeper/testdata/test-contract
	$(MAKE) -C ./x/compute/internal/keeper/testdata/ibc-test-contract
	cp ./x/compute/internal/keeper/testdata/ibc-test-contract/ibc.wasm ./x/compute/internal/keeper/testdata/
	cat ./x/compute/internal/keeper/testdata/contract.wasm | gzip -9 > ./x/compute/internal/keeper/testdata/contract.wasm.gz

prep-go-tests: build-test-contracts  bin-data-sw
	# empty BUILD_PROFILE means debug mode which compiles faster
	SGX_MODE=SW $(MAKE) build-linux
	cp ./cosmwasm/enclaves/execute/librust_cosmwasm_enclave.signed.so ./x/compute/internal/keeper

go-tests: build-test-contracts bin-data-sw
	# empty BUILD_PROFILE means debug mode which compiles faster
	SGX_MODE=SW $(MAKE) build-linux
	cp ./cosmwasm/enclaves/execute/librust_cosmwasm_enclave.signed.so ./x/compute/internal/keeper
	rm -rf ./x/compute/internal/keeper/trustlesshub/.sgx_secrets
	mkdir -p ./x/compute/internal/keeper/trustlesshub/.sgx_secrets
	GOMAXPROCS=8 SGX_MODE=SW TRST_SGX_STORAGE='./' go test -failfast -timeout 60m -v ./x/compute/internal/... -tags test

go-tests-hw: build-test-contracts bin-data
	# empty BUILD_PROFILE means debug mode which compiles faster
	SGX_MODE=HW $(MAKE) build-linux
	cp ./cosmwasm/enclaves/execute/librust_cosmwasm_enclave.signed.so ./x/compute/internal/keeper
	rm -rf ./x/compute/internal/keeper/trustlesshub/.sgx_secrets
	mkdir -p ./x/compute/internal/keeper/trustlesshub/.sgx_secrets
	GOMAXPROCS=8 SGX_MODE=HW go test -v ./x/compute/internal/... $(GO_TEST_ARGS)

# When running this more than once, after the first time you'll want to remove the contents of the `ffi-types`
# rule in the Makefile in `enclaves/execute`. This is to speed up the compilation time of tests and speed up the
# test debugging process in general.
.PHONY: enclave-tests
enclave-tests:
	$(MAKE) -C cosmwasm/enclaves/test run

build-all-test-contracts: build-test-contracts
	cd ./cosmwasm/contracts/hackatom && RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown
	wasm-opt -Os ./cosmwasm/contracts/hackatom/target/wasm32-unknown-unknown/release/hackatom.wasm -o ./x/compute/internal/keeper/testdata/hackatom.wasm

	cd ./cosmwasm/contracts/staking && RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown 
	wasm-opt -Os ./cosmwasm/contracts/staking/target/wasm32-unknown-unknown/release/staking.wasm -o ./x/compute/internal/keeper/testdata/staking.wasm


build-non-test-contracts: build-test-contracts
	cd ./cosmwasm/contracts/ibc-reflect && RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown 
	wasm-opt -Os ./cosmwasm/contracts/reflect/target/wasm32-unknown-unknown/release/reflect.wasm -o ./x/compute/internal/keeper/testdata/ibc-reflect.wasm

	cd ./cosmwasm/contracts/burner && RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown 
	wasm-opt -Os ./cosmwasm/contracts/burner/target/wasm32-unknown-unknown/release/burner.wasm -o ./x/compute/internal/keeper/testdata/burner.wasm

bin-data: bin-data-sw bin-data-develop bin-data-production

bin-data-sw:
	cd ./cmd/trstd && go-bindata -o ias_bin_sw.go -prefix "../../ias_keys/sw_dummy/" -tags "!hw" ../../ias_keys/sw_dummy/...

bin-data-develop:
	cd ./cmd/trstd && go-bindata -o ias_bin_dev.go -prefix "../../ias_keys/develop/" -tags "develop,hw" ../../ias_keys/develop/...

bin-data-production:
	cd ./cmd/trstd && go-bindata -o ias_bin_prod.go -prefix "../../ias_keys/production/" -tags "production,hw" ../../ias_keys/production/...
	
TAG := 0.12.0
trustless-contract-optimizer:
	docker build -f deployment/dockerfiles/trustless-contract-optimizer.Dockerfile -t trstlabs/trustless-contract-optimizer:${TAG} .
	docker tag trstlabs/trustless-contract-optimizer:${TAG} trstlabs/trustless-contract-optimizer:latest

aesm-image:
	docker build -f deployment/dockerfiles/aesm.Dockerfile -t trstlabs/aesm .

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
	rly tx connection trstdev1-trstdev2 --override
	rly start trstdev1-trstdev2 -p events -b 100 --debug > rly.log

kill-dev:
	@echo "Killing trstd and removing previous data"
	-@rm -rf ./data
	-@killall trstd 2>/dev/null

# starts Cosmos relayer given the localchains are running
run-rly: 
	@echo "Starting up local test relayer..."
	./deployment/ibc/start.sh

# starts a Cosmos relayer, localtrst 1&2 and localjuno
run-rly-juno: 
	@echo "Starting up local test relayers..."
	docker compose -f deployment/ibc/docker-compose-rly.yml up

# stops  a Cosmos relayer, localtrst 1&2 and localjuno
kill-rly-juno: 
	@echo "Starting up local test relayers..."
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
###                         Swagger & Protobuf                              ###
###############################################################################

.PHONY: update-swagger-openapi-docs statik statik-install proto-swagger-openapi-gen

statik-install:
	@echo "Installing statik..."
	@go install github.com/rakyll/statik@v0.1.6

statik:
	statik -src=client/docs/static/ -dest=client/docs -f -m

proto-swagger-openapi-gen:
	cp go.mod /tmp/go.mod.bak
	cp go.sum /tmp/go.sum.bak
	@./scripts/protoc-swagger-openapi-gen.sh
	cp /tmp/go.mod.bak go.mod
	cp /tmp/go.sum.bak go.sum

# Example `CHAIN_VERSION=v1.4.0 make update-swagger-openapi-docs`
update-swagger-openapi-docs: statik-install proto-swagger-openapi-gen statik

protoVer=v0.2

proto-all: proto-lint proto-gen proto-swagger-openapi-gen

proto-gen:
	cp go.mod /tmp/go.mod.bak
	cp go.sum /tmp/go.sum.bak
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen:$(protoVer) sh ./scripts/protocgen.sh
	cp /tmp/go.mod.bak go.mod
	cp /tmp/go.sum.bak go.sum
	go mod tidy

proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

.PHONY: proto-all proto-gen proto-format proto-lint proto-check-breaking