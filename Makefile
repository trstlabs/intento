PACKAGES=$(shell go list ./... | grep -v '/simulation')
VERSION ?= $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
CURRENT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
LEDGER_ENABLED ?= false
BINDIR ?= $(GOPATH)/bin
BUILD_PROFILE ?= release
DEB_BIN_DIR ?= /usr/local/bin
DEB_LIB_DIR ?= /usr/lib

SGX_MODE ?= HW
BRANCH ?= develop
DEBUG ?= 0


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

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=tpp \
-X github.com/cosmos/cosmos-sdk/version.AppName=tppd \
	-X github.com/danieljdd/tpp/cmd/tppd/version.ClientName=tppd \
-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \


ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ldflags += -s -w
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

GO_TAGS := $(build_tags)
# -ldflags
LD_FLAGS := $(ldflags)

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install

install: go.sum
	@echo "--> Installing tppd"
	@go install $(BUILD_FLAGS) ./cmd/tppd

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

test:
	@go test -mod=readonly $(PACKAGES)

vendor:
	cargo vendor third_party/vendor --manifest-path third_party/build/Cargo.toml

build_local_no_rust: bin-data-$(IAS_BUILD)
	cp go-cosmwasm/target/release/libgo_cosmwasm.so go-cosmwasm/api
	go build -mod=readonly -tags "$(GO_TAGS)" -ldflags '$(LD_FLAGS)' ./cmd/tppd

build-linux: vendor bin-data-$(IAS_BUILD)
	BUILD_PROFILE=$(BUILD_PROFILE) $(MAKE) -C go-cosmwasm build-rust
	cp go-cosmwasm/target/$(BUILD_PROFILE)/libgo_cosmwasm.so go-cosmwasm/api
#   this pulls out ELF symbols, 80% size reduction!
	go build -mod=readonly -tags "$(GO_TAGS)" -ldflags '$(LD_FLAGS)' ./cmd/tppd


deb: build-linux deb-no-compile

deb-no-compile:

	rm -rf /tmp/tpp

	mkdir -p /tmp/tpp/deb/$(DEB_BIN_DIR)
	mv -f ./tppd /tmp/tpp/deb/$(DEB_BIN_DIR)/tppd
	mkdir -p /tmp/tpp/deb/$(DEB_LIB_DIR)
	cp -f ./go-cosmwasm/api/libgo_cosmwasm.so ./go-cosmwasm/librust_cosmwasm_enclave.signed.so /tmp/tpp/deb/$(DEB_LIB_DIR)/
	chmod +x /tmp/tpp/deb/$(DEB_LIB_DIR)/lib*.so

	mkdir -p /tmp/tpp/deb/DEBIAN
	cp ./deployment/deb/control /tmp/tpp/deb/DEBIAN/control
	printf "Version: " >> /tmp/tpp/deb/DEBIAN/control
	printf "$(VERSION)" >> /tmp/tpp/deb/DEBIAN/control
	echo "" >> /tmp/tpp/deb/DEBIAN/control
	cp ./deployment/deb/postinst /tmp/tpp/deb/DEBIAN/postinst
	chmod 755 /tmp/tpp/deb/DEBIAN/postinst
	cp ./deployment/deb/postrm /tmp/tpp/deb/DEBIAN/postrm
	chmod 755 /tmp/tpp/deb/DEBIAN/postrm
	cp ./deployment/deb/triggers /tmp/tpp/deb/DEBIAN/triggers
	chmod 755 /tmp/tpp/deb/DEBIAN/triggers
	dpkg-deb --build /tmp/tpp/deb/ .
	-rm -rf /tmp/tpp

rename_for_release:
	-rename "s/windows-4.0-amd64/v${VERSION}-win64/" *.exe
	-rename "s/darwin-10.6-amd64/v${VERSION}-osx64/" *darwin*

sign_for_release: rename_for_release
	sha256sum tpp-blockchain*.deb > SHA256SUMS
	-sha256sum tppd-* >> SHA256SUMS
	gpg -u 91831DE812C6415123AFAA7B420BF1CB005FBCE6 --digest-algo sha256 --clearsign --yes SHA256SUMS
	rm -f SHA256SUMS

release: sign_for_release
	rm -rf ./release/
	mkdir -p ./release/
	cp tpp-blockchain_*.deb ./release/
	cp tppd-* ./release/
	cp SHA256SUMS.asc ./release/

clean:
	-rm -rf /tpp

	-rm -f ./tppd*
#   -find -name librust_cosmwasm_enclave.signed.so -delete
#   -find -name libgo_cosmwasm.so -delete
#   -find -name '*.so' -delete
#   -find -name 'target' -type d -exec rm -rf \;
	-rm -f ./tpp-blockchain*.deb
	-rm -f ./SHA256SUMS*
	-rm -rf ./third_party/vendor/
	-rm -rf ./.sgx_secrets/*
	-rm -rf ./x/compute/internal/keeper/.sgx_secrets/*
	-rm -rf ./*.der
	-rm -rf ./x/compute/internal/keeper/*.der
	-rm -rf ./cmd/tppd/ias_bin*
	$(MAKE) -C go-cosmwasm clean-all
	$(MAKE) -C cosmwasm/packages/wasmi-runtime clean

# while developing:
build-enclave: vendor
	$(MAKE) -C cosmwasm/packages/wasmi-runtime

# while developing:
check-enclave:
	$(MAKE) -C cosmwasm/packages/wasmi-runtime check

# while developing:
clippy-enclave:
	$(MAKE) -C cosmwasm/packages/wasmi-runtime clippy

# while developing:
clean-enclave:
	$(MAKE) -C cosmwasm/packages/wasmi-runtime clean

sanity-test:
	SGX_MODE=SW $(MAKE) build-linux
	cp ./cosmwasm/packages/wasmi-runtime/librust_cosmwasm_enclave.signed.so .
	SGX_MODE=SW ./cosmwasm/testing/sanity-test.sh

sanity-test-hw:
	$(MAKE) build-linux
	cp ./cosmwasm/packages/wasmi-runtime/librust_cosmwasm_enclave.signed.so .
	./cosmwasm/testing/sanity-test.sh

callback-sanity-test:
	SGX_MODE=SW $(MAKE) build-linux
	cp ./cosmwasm/packages/wasmi-runtime/librust_cosmwasm_enclave.signed.so .
	SGX_MODE=SW ./cosmwasm/testing/callback-test.sh

build-test-contract:
# echo "" | sudo add-apt-repository ppa:hnakamur/binaryen
# sudo apt update
# sudo apt install -y binaryen
	$(MAKE) -C ./x/compute/internal/keeper/testdata/test-contract

prep-go-tests: build-test-contract
# empty BUILD_PROFILE means debug mode which compiles faster
	SGX_MODE=SW $(MAKE) build-linux
	cp ./cosmwasm/packages/wasmi-runtime/librust_cosmwasm_enclave.signed.so ./x/compute/internal/keeper

go-tests: build-test-contract
# empty BUILD_PROFILE means debug mode which compiles faster
	SGX_MODE=SW $(MAKE) build-linux
	cp ./cosmwasm/packages/wasmi-runtime/librust_cosmwasm_enclave.signed.so ./x/compute/internal/keeper
	rm -rf ./x/compute/internal/keeper/.sgx_secrets
	mkdir -p ./x/compute/internal/keeper/.sgx_secrets
	SGX_MODE=SW go test -timeout 1200s -p 1 -v ./x/compute/internal/... $(GO_TEST_ARGS)

go-tests-hw: build-test-contract
# empty BUILD_PROFILE means debug mode which compiles faster
	SGX_MODE=HW $(MAKE) build-linux
	cp ./cosmwasm/packages/wasmi-runtime/librust_cosmwasm_enclave.signed.so ./x/compute/internal/keeper
	rm -rf ./x/compute/internal/keeper/.sgx_secrets
	mkdir -p ./x/compute/internal/keeper/.sgx_secrets
	SGX_MODE=HW go test -p 1 -v ./x/compute/internal/... $(GO_TEST_ARGS)

	.PHONY: enclave-tests
enclave-tests:
	$(MAKE) -C cosmwasm/packages/enclave-test run

build-all-test-contracts: build-test-contract
# echo "" | sudo add-apt-repository ppa:hnakamur/binaryen
# sudo apt update
# sudo apt install -y binaryen
	cd ./cosmwasm/contracts/gov && RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown --locked
	wasm-opt -Os ./cosmwasm/contracts/gov/target/wasm32-unknown-unknown/release/gov.wasm -o ./x/compute/internal/keeper/testdata/gov.wasm

	cd ./cosmwasm/contracts/dist && RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown --locked
	wasm-opt -Os ./cosmwasm/contracts/dist/target/wasm32-unknown-unknown/release/dist.wasm -o ./x/compute/internal/keeper/testdata/dist.wasm

	cd ./cosmwasm/contracts/mint && RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown --locked
	wasm-opt -Os ./cosmwasm/contracts/mint/target/wasm32-unknown-unknown/release/mint.wasm -o ./x/compute/internal/keeper/testdata/mint.wasm

	cd ./cosmwasm/contracts/staking && RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown --locked
	wasm-opt -Os ./cosmwasm/contracts/staking/target/wasm32-unknown-unknown/release/staking.wasm -o ./x/compute/internal/keeper/testdata/staking.wasm

	cd ./cosmwasm/contracts/reflect && RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown --locked
	wasm-opt -Os ./cosmwasm/contracts/reflect/target/wasm32-unknown-unknown/release/reflect.wasm -o ./x/compute/internal/keeper/testdata/reflect.wasm

	cd ./cosmwasm/contracts/burner && RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown --locked
	wasm-opt -Os ./cosmwasm/contracts/burner/target/wasm32-unknown-unknown/release/burner.wasm -o ./x/compute/internal/keeper/testdata/burner.wasm

	cd ./cosmwasm/contracts/erc20 && RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown --locked
	wasm-opt -Os ./cosmwasm/contracts/erc20/target/wasm32-unknown-unknown/release/cw_erc20.wasm -o ./x/compute/internal/keeper/testdata/erc20.wasm

	cd ./cosmwasm/contracts/hackatom && RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown --locked
	wasm-opt -Os ./cosmwasm/contracts/hackatom/target/wasm32-unknown-unknown/release/hackatom.wasm -o ./x/compute/internal/keeper/testdata/contract.wasm
	cat ./x/compute/internal/keeper/testdata/contract.wasm | gzip > ./x/compute/internal/keeper/testdata/contract.wasm.gzip

bin-data: bin-data-sw bin-data-develop bin-data-production

bin-data-sw:
	cd ./cmd/tppd && go-bindata -o ias_bin_sw.go -prefix "../../ias_keys/sw_dummy/" -tags "!hw" ../../ias_keys/sw_dummy/...

bin-data-develop:
	cd ./cmd/tppd && go-bindata -o ias_bin_dev.go -prefix "../../ias_keys/develop/" -tags "develop,hw" ../../ias_keys/develop/...

bin-data-production:
	cd ./cmd/tppd && go-bindata -o ias_bin_prod.go -prefix "../../ias_keys/production/" -tags "production,hw" ../../ias_keys/production
