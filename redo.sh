sudo apt remove trustless-hub-node -y
make clean-files
sudo apt install ./trustless-hub-node_0.7.7-69-g4a408b4_amd64.deb -y

./inittest.sh