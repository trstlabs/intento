sudo apt remove trustless-hub-node -y
make clean-files
sudo apt install ./trustless-hub-node_0.7.5-22-g70baf50_amd64.deb -y

./init.sh
