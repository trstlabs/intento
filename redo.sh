sudo apt remove trustless-hub-node -y
make clean-files
sudo apt install ./trustless-hub-node_0.7.7-23-g1536ae8_amd64.deb -y

./inittest.sh
