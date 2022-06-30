sudo apt remove trustless-hub-node -y
make clean-files
sudo apt install ./trustless-hub-node_0.7.7-7-g8ef7c64_amd64.deb -y

./inittest.sh
