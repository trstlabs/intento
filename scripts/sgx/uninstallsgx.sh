sudo apt purge -y libsgx-enclave-common libsgx-enclave-common-dev libsgx-urts sgx-aesm-service libsgx-uae-service libsgx-launch libsgx-aesm-launch-plugin libsgx-ae-le
sudo /opt/intel/sgxdriver/uninstall.sh
sudo "$HOME"/.sgxsdk/sgxsdk/uninstall.sh
rm -rf "$HOME/.sgxsdk"
sudo apt purge -y libsgx-enclave-common libsgx-enclave-common-dev libsgx-urts sgx-aesm-service libsgx-uae-service libsgx-launch libsgx-aesm-launch-plugin libsgx-ae-le
rm -rf "$HOME/var/opt/aesmd"
sudo apt autoremove
