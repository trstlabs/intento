sudo apt remove trustless-hub-node -y
make clean-files
sudo apt install ./trustless-hub-node_0.7.7-25-gc9f88b0_amd64.deb -y

./inittest.sh
