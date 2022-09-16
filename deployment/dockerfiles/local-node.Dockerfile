# Final image
FROM build-release

ARG SGX_MODE=SW
ENV SGX_MODE=${SGX_MODE}
#
ARG SECRET_LOCAL_NODE_TYPE
ENV SECRET_LOCAL_NODE_TYPE=${SECRET_LOCAL_NODE_TYPE}

ENV PKG_CONFIG_PATH=""
ENV TRST_ENCLAVE_DIR=/usr/lib/

COPY deployment/docker/sanity-test.sh /root/
RUN chmod +x /root/sanity-test.sh

COPY x/compute/internal/keeper/testdata/erc20.wasm erc20.wasm
RUN true
COPY deployment/ci/wasmi-sgx-test.sh .
RUN true
COPY deployment/ci/bootstrap_init.sh .
RUN true
COPY deployment/ci/node_init.sh .
RUN true
COPY deployment/ci/startup.sh .
RUN true
COPY deployment/ci/node_key.json .

RUN chmod +x /usr/bin/trstd
# RUN chmod +x /usr/bin/trstcli
RUN chmod +x wasmi-sgx-test.sh
RUN chmod +x bootstrap_init.sh
RUN chmod +x startup.sh
RUN chmod +x node_init.sh


#RUN mkdir -p /root/.trstd/.compute/
#RUN mkdir -p /root/.sgx_secrets/
#RUN mkdir -p /root/.trstd/.node/

# Enable autocomplete
#RUN trstcli completion > /root/trstcli_completion
#RUN trstd completion > /root/trstd_completion
#
#RUN echo 'source /root/trstd_completion' >> ~/.bashrc
#RUN echo 'source /root/trstcli_completion' >> ~/.bashrc

#ENV LD_LIBRARY_PATH=/opt/sgxsdk/libsgx-enclave-common/:/opt/sgxsdk/lib64/

# Run trstd by default, omit entrypoint to ease using container with trstcli
ENTRYPOINT ["/bin/bash", "startup.sh"]
