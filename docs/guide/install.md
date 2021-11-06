---
order: 1
description: Steps to install a Trustless Hub Daemon on your local computer.
---

[Install Trustless Hub](#install-trustless-hub)
  - [Prerequisites](#prerequisites)
    - [Operating Systems](#operating-systems)
    - [Go and Rust](#go-and-rust)
    - [SGX tools](#sgx-tools)
    - [Using `sgx-detect`](Using-`sgx-detect`)

  - [Install Trustless Hub daemon](#install-trustless-hub-daemon)
    - [Install the debian package](#Install-the-debian-package)
  - [Uninstall](#uninstall-trustless-hub)



# Install Trustless Hub

You can install the Trustless Hub (https://github.com/trstlabs/trst) blockchain on a server or you can install Trustless Hub on your local computer.  


## Prerequisites

Be sure you have met the prerequisites before you install and use Trustless Hub. 

::: danger The Trustless Hub daemon requires an Intel SGX enabled device. This has been tested on a dedicated server running Linux 20.04.3 LTS. 
:::

1. Go to your BIOS 
2. Enable SGX
3. Disable Secure Boot
4. Disbale HyperV if available


### Operating System

Trustless Hub is supported for the following operating systems:

- GNU/Linux


## Go and Rust

Trustless Hub is written in the Go and Rust programming language. To use Trustless Hub on a local system:

- Install [Go](https://golang.org/doc/install) (**version 1.16** or higher)
- Ensure the Go environment variables are [set properly](https://golang.org/doc/gopath_code#GOPATH) on your system


- Install [Rustup](https://www.rust-lang.org/tools/install) (**version 1.16** or higher)
First, make sure you have Rust installed: https://www.rust-lang.org/tools/install

Once Rust is installed, install the nightly toolchain:

rustup toolchain install nightly

```sh
rustup install nightly-2020-10-25
rustup default nightly-2020-10-25
rustup component add rust-src
cargo +1.49.0-nightly run

rustup override set nightly-2020-10-25
```

## SGX tools

Install the latest version of the SGX tools by running the following script. Be sure to check if it contains the latest [Intel SGX version](https://download.01.org/intel-sgx/sgx-linux/) for security reasons.  Install:

```sh
bash "$HOME/trst/installsgx.sh"
```


## Using `sgx-detect`:

- Once Rust is installed, install the `nightly` toolchain:

```bash
rustup toolchain install nightly-2020-10-25
```
Then install the SGX tools:
```bash
sudo apt install -y libssl-dev protobuf-compiler
cargo +nightly install fortanix-sgx-tools sgxs-tools

sgx-detect
```

When you do run sgx-detect, it should print at the end:
AES Service should be running and all checks should be green. If this is not the case, try reinstalling using a different script. 

```
✔  Able to launch enclaves
   ✔  Debug mode
   ✔  Production mode (Intel whitelisted)

You're all set to start running SGX programs!
```


# Install Trustless Hub daemon

First, download the Trustless Hub daemon, run the following command in the home directory:

```sh
curl https://github.com/trstlabs/trst! | bash
cd trst

```


## Install the debian package

To install Trustless Hub daemon, run the following command in the trst directoty:

```sh
make deb
sudo apt install ./trustlesshub_[version].deb -y

```

## Verify Your Trustless Hub Deaemon Version 

To verify the version of Trustless Hub you have installed, run the following command:

```sh
trstd version
```

Success! You are now ready to start a local node or join the mainnet! :)

# Uninstall Trustless Hub

To uninstall Trustless Hub, run the following command in the trst directoty:

```sh
sudo apt remove trustlesshub -y
make clean
```

# Uninstall SGX

To uninstall all SGX related tools, run:
```sh
bash "$HOME/trst/unnstallsgx.sh"
```

Or else, to uninstall the Intel(R) SGX Driver manually, run:

```bash
sudo /opt/intel/sgxdriver/uninstall.sh
```

The above command produces no output when it succeeds. If you want to verify that the driver has been uninstalled, you can run the following, which should print `SGX Driver NOT installed`:

```bash
ls /dev/isgx &>/dev/null && echo "SGX Driver installed" || echo "SGX Driver NOT installed"
```

To uninstall the SGX SDK, run:

```bash
sudo "$HOME"/.sgxsdk/sgxsdk/uninstall.sh
rm -rf "$HOME/.sgxsdk"
```

To uninstall the rest of the dependencies, run:

```bash
sudo apt purge -y libsgx-enclave-common libsgx-enclave-common-dev libsgx-urts sgx-aesm-service libsgx-uae-service libsgx-launch libsgx-aesm-launch-plugin libsgx-ae-le
```


See the [Secret Network docs](https://build.scrt.network/validators-and-full-nodes/setup-sgx.html#for-contract-developers) for additional information related to Rust and SGX.

# Install
