sudo apt remove trustless-hub-node -y
make clean-files
sudo apt install ./trustless-hub-node_0.7.5-16-g00eaaf2_amd64.deb -y

./inittest.sh
