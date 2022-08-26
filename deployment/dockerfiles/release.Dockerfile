# Base image
FROM rust-go-base-image AS build-env-rust-go

# Final image
FROM trstlabs/sgx-base-trustless_hub:2004-1.1.3 as build-release

# wasmi-sgx-test script requirements
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    #### Base utilities ####
    jq \
    openssl \
    curl \
    wget \
    libsnappy-dev \
    libgflags-dev \
    bash-completion && \
    rm -rf /var/lib/apt/lists/*

RUN echo "source /etc/profile.d/bash_completion.sh" >> ~/.bashrc

RUN curl -sL https://deb.nodesource.com/setup_15.x | bash - && \
    apt-get update && \
    apt-get install -y nodejs npm && \
    npm i -g local-cors-proxy

ARG SGX_MODE=SW
ENV SGX_MODE=${SGX_MODE}

ARG TRST_NODE_TYPE=BOOTSTRAP
ENV TRST_NODE_TYPE=${TRST_NODE_TYPE}

ENV TRST_ENCLAVE_DIR=/usr/lib/

# workaround because paths seem kind of messed up
RUN cp /opt/sgxsdk/lib64/libsgx_urts_sim.so /usr/lib/libsgx_urts_sim.so
RUN cp /opt/sgxsdk/lib64/libsgx_uae_service_sim.so /usr/lib/libsgx_uae_service_sim.so

# Install ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env-rust-go /go/src/github.com/trstlabs/SecretNetwork/go-cosmwasm/target/release/libgo_cosmwasm.so /usr/lib/
COPY --from=build-env-rust-go /go/src/github.com/trstlabs/SecretNetwork/go-cosmwasm/librust_cosmwasm_enclave.signed.so /usr/lib/
COPY --from=build-env-rust-go /go/src/github.com/trstlabs/SecretNetwork/go-cosmwasm/librust_cosmwasm_query_enclave.signed.so /usr/lib/
COPY --from=build-env-rust-go /go/src/github.com/trstlabs/SecretNetwork/trstd /usr/bin/trstd

COPY deployment/docker/bootstrap/bootstrap_init.sh .
COPY deployment/docker/node/node_init.sh .
COPY deployment/docker/startup.sh .
COPY deployment/docker/node_key.json .

RUN chmod +x /usr/bin/trstd
RUN chmod +x bootstrap_init.sh
RUN chmod +x startup.sh
RUN chmod +x node_init.sh

RUN trstd completion > /root/trstd_completion

RUN echo 'source /root/trstd_completion' >> ~/.bashrc

RUN mkdir -p /root/.trstd/.compute/
RUN mkdir -p /opt/trustlesshub/.sgx_secrets/
RUN mkdir -p /root/.trstd/.node/
RUN mkdir -p /root/config/



####### Node parameters
ARG MONIKER=default
ARG CHAINID=trst_chain_1
ARG GENESISPATH=https://raw.githubusercontent.com/trstlabs/SecretNetwork/master/secret-testnet-genesis.json
ARG PERSISTENT_PEERS=201cff36d13c6352acfc4a373b60e83211cd3102@bootstrap.southuk.azure.com:26656

ENV GENESISPATH="${GENESISPATH}"
ENV CHAINID="${CHAINID}"
ENV MONIKER="${MONIKER}"
ENV PERSISTENT_PEERS="${PERSISTENT_PEERS}"

#ENV LD_LIBRARY_PATH=/opt/sgxsdk/libsgx-enclave-common/:/opt/sgxsdk/lib64/

# Run trstd by default, omit entrypoint to ease using container with trstcli
ENTRYPOINT ["/bin/bash", "startup.sh"]
