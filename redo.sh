sudo apt remove trustless-hub-node -y
make clean-files
sudo apt install ./trustless-hub-node_0.7.7-95-g42545324_amd64.deb -y

./inittest.sh