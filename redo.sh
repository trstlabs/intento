sudo apt remove trustless-hub-node -y
make clean-files
sudo apt install ./trustless-hub-node_0.7.7-12-g22ffe06_amd64.deb -y

./inittest.sh
