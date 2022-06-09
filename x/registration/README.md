<!--
order: 0
-->

# Registration module
The registration module contains code related to the remote attestation of a new secure enclave instance to the instances of the enclave in in the chain. It is used to securely share secrets between enclaves.


# Enclave keys
Several keys are generated in the enclave for different purposes.
From the bootstrap node's enclave private key, the consensus seed is derived. From the consensus seed, 4 types of keys are generated that are used within the network for deriving, attesting or sharing various types of data between nodes and users. 
Consensus_seed_exchange_keypair, for exchanging the consensus seed through the verifying node’s pubkey, and then validated by its privkey along with a valid certificate.

1. Consensus_io_exchange_keypair, to encrypt transaction inputs, with the sender pubkey, then decrypt these transaction inputs through the privkey and the sender signature and encrypt transaction outputs again with the sender pubkey. 
2. Consensus_state_ikm, to help encrypt and decrypt the contract state. This is done by combining this with a unique contract key.
3. Consensus_callback_secret, to create callback signatures with. The callback signature is a sha236 hash with: this key, extended with the to-execute address,the message to sign as well as the funds that have to be sent. The message to sign is an encrypted message (through the Consensus_io_exchange_pubkey + nonce + sender pubkey) and prepended with a hash of the code. The callback signature is verified by recreating it in the enclave with the encrypted input message, sender address, and the funds. 

Each contract has a unique and unforgettable Contract_key. Contract_key is a concatenation of two values: signer_id and authenticated_contract_key. Each contract has a unique unforgeable encryption key. It is unique as the state of two contracts with the same code is different. It is unforgettable as a malicious node runner won't be able to locally encrypt transactions with it's own encryption key, and then decrypt the newly created state with a malicious key.

# Bootstrapping the first node with the enclave
## InitBootstrapCmd 
generates the master private/public key as well as 4 other keys. This happens once at the initialization of a chain. It requires a spid and api key. It returns the consensus_seed_exchange_keypair public key, which is saved on-chain, and used to propagate the seed to registering nodes.

It performs 3 functions:
Create_consensus_seed, which creates the consensus seed
Generate_consensus_master_keys, creates all of the keypairs (consensus_seed_exchange_keypair,consensus_io_exchange_keypair, consensus_state_ikm, consensus_callback_secret).
Create_registration_key -> Registration KeyPair
These 3 functions are the same for a new node joining the chain.
Thereafter, it creates an attestation report, using the valid SPID and API phrases provided.
Hereafter the consensus_io_exchange_keypair public key is created and returned, the private key stays in the enclave and the public key is revealed and written on-chain.

# Registration of new node and enclave
### InitAttestationCmd 
creates an attestation report that verifies the enclave’s integrity. It requires a spid and api key. In the enclave, it creates a keypair through SgxEccHandle from the apache teaclave-sgx-sdk rust library. Intel signes a report through EPID and the API key provided. Following this report, we get a quote from the quoting enclave. These exact steps are found in cosmwasm/enclaves/execute/src.registration/attestation.rs.
“The Quoting Enclave creates the EPID key used for signing platform attestations which is then certified by an EPID backend infrastructure. The EPID key represents not only the platform but the trustworthiness of the underlying hardware. Only the Quoting Enclave has access to the EPID key when the enclave system is operational, and the EPID key is bound to the version of the processor’s firmware. Therefore, a QUOTE can be seen to be issued by the processor itself.” (Intel via https://www.intel.com/content/www/us/en/developer/articles/technical/innovative-technology-for-cpu-based-attestation-and-sealing.html)
The output of this function is an X.509 certificate signed by the enclave, which contains the report signed by Intel, it is written to opt/trustlesshub/.sgx_secrets by default.


### AuthenticateNodeCmd
uploads the previously generated EPID-specific certificate on-chain to authenticate the node. All nodes verify if the signature is valid and then will try to register the node’s enclave. If the node is already registered it will get the seed on-chain. Else, the existing node will authenticate the new node through ecall_authenticate_new_node. The in-enclave function will authenticate the new node based on a received certificate. The consensus seed will be encrypted with the new node’s public key and shared on-chain if the node’s certificate is authenticated successfully.

### ConfigureCredentialsCmd 
configures the node with the credentials file and the encrypted seed that was written on-chain. It writes to $HOME/.trst/.node/seed.json by default.

# Init 
opon node initiation, the enclave will take the seed, stored at $HOME/.trst/.node/seed.json. In the enclave, init_node is called with the certificate and the seed, it sets the consensus seed by decrypting it from the enclave keypair after verifying the certificate again. From this it also creates all of the keypairs (consensus_seed_exchange_keypair,consensus_io_exchange_keypair, consensus_state_ikm, consensus_callback_secret).
