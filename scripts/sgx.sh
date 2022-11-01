#! /bin/bash

UBUNTUVERSION=$(lsb_release -r -s | cut -d '.' -f 1)

if (($UBUNTUVERSION < 16)); then
	echo "Your version of Ubuntu is not supported. Must have Ubuntu 16.04 and up. Aborting installation script..."
	exit 1
elif (($UBUNTUVERSION < 18)); then
	DISTRO='xenial'
else
	DISTRO='bionic'
fi

echo "\n\n#######################################"
echo "##### Installing missing packages #####"
echo "#######################################\n\n"

# Install needed packages for script
sudo apt install -y lynx parallel gdebi make

# Create a working directory to download and install the SDK inside
mkdir -p "$HOME/.sgxsdk"

(
   # In a new sub-shell cd into our working directory so to no pollute the
   # original shell's working directory
   cd "$HOME/.sgxsdk"

   echo "\n\n################################################"
   echo "##### Downloading Intel SGX driver and SDK #####"
   echo "################################################\n\n"

   # 1. Go to https://download.01.org/intel-sgx/sgx-linux
   # 2. Step into the latest version
   # 3. Step into `distro/$LATEST_UBUNTU_YOU_SEE_THERE`
   # 4. Download `sgx_linux_x64_driver_*.bin` and `sgx_linux_x64_sdk_*.bin`
   lynx -dump -listonly -nonumbers https://download.01.org/intel-sgx/sgx-linux/ |
      grep -P 'sgx-linux/(\d\.?)+/' |
      sort -V |
      tail -1 |
      parallel --bar --verbose lynx -dump -listonly -nonumbers "{}/distro" |
      grep -P 'ubuntu\d\d' |
      sort -V |
      tail -1 |
      parallel --bar --verbose lynx -dump -listonly -nonumbers |
      grep -P '\.bin$' |
      parallel --bar --verbose curl -OSs

   # Make the driver and SDK installers executable
   chmod +x ./sgx_linux_*.bin

   echo "\n\n###############################################"
   echo "##### Installing Intel SGX driver and SDK #####"
   echo "###############################################\n\n"

   # Install the driver
   sudo ./sgx_linux_x64_driver_*.bin

   # Remount /dev as exec, also at system startup
   sudo tee /etc/systemd/system/remount-dev-exec.service >/dev/null <<EOF
[Unit]
Description=Remount /dev as exec to allow AESM service to boot and load enclaves into SGX

[Service]
Type=oneshot
ExecStart=/bin/mount -o remount,exec /dev
RemainAfterExit=true

[Install]
WantedBy=multi-user.target
EOF
   sudo systemctl enable remount-dev-exec
   sudo systemctl start remount-dev-exec

   # Install the SDK inside ./sgxsdk/ which is inside $HOME/.sgxsdk
   echo yes | ./sgx_linux_x64_sdk_*.bin

   # Setup the environment variables for every new shell
   echo "source '$HOME/.sgxsdk/sgxsdk/environment'" |
      tee -a "$HOME/.bashrc" "$HOME/.zshrc" > /dev/null
)

echo "\n\n##############################################"
echo "##### Installing additional dependencies #####"
echo "##############################################\n\n"

# Add Intels's SGX PPA
echo "deb [arch=amd64] https://download.01.org/intel-sgx/sgx_repo/ubuntu $DISTRO main" |
   sudo tee /etc/apt/sources.list.d/intel-sgx.list
wget -qO - https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key |
   sudo apt-key add -
sudo apt update

# Install all the additional necessary dependencies (besides the driver and the SDK)
# for building a rust enclave
wget -O /tmp/libprotobuf10_3.0.0-9_amd64.deb http://ftp.br.debian.org/debian/pool/main/p/protobuf/libprotobuf10_3.0.0-9_amd64.deb
(sleep 3 ; echo y) | sudo gdebi /tmp/libprotobuf10_3.0.0-9_amd64.deb

sudo apt install -y libsgx-enclave-common libsgx-enclave-common-dev libsgx-urts sgx-aesm-service libsgx-uae-service libsgx-launch libsgx-aesm-launch-plugin libsgx-ae-le autoconf libtool